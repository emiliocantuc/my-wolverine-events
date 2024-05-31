package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Database connection
var db *DB

// TODO do we need this?
func respondWithTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		t.Execute(w, data)
	}
}

func login(w http.ResponseWriter, req *http.Request) {

	handleError := func(err error) bool {
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "Could not authenticate")
			return true
		}
		return false
	}

	// Verify the Signin with google token
	email, err := verifySignin(req.FormValue("credential"))
	if handleError(err) {
		return
	}

	// Check if user exsits
	exists, err := db.UserExists(email)
	if handleError(err) {
		return
	}

	// Create user email not registered
	var id int64
	redirectTo := "/"
	if !exists {
		redirectTo = "/prefs"
		id, err = db.AddUser(email)
		fmt.Println("New user", email)
		if handleError(err) {
			return
		}
	}

	// Generate JWT of email, place in cookie "jwt"
	jwtString, err := jwtEncodeEmailId(email, id)
	if handleError(err) {
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: jwtString})
	http.Redirect(w, req, redirectTo, http.StatusSeeOther)
}

func logout(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: ""})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, req *http.Request) {

	var featuredEvents, recommendedEvents []EventCard
	var err error

	handleError := func(err error) bool {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal error")
			return true
		}
		return false
	}

	setLoggedIn := func(events []EventCard, loggedIn bool) {
		for i := 0; i < len(events); i++ {
			events[i].LoggedIn = loggedIn
		}
	}

	loggedIn, email, id := loggedInEmail(req)
	fmt.Println(loggedIn, email, id)

	if loggedIn {
		recommendedEvents, err = db.GetRecommendedEvents(id)
		fmt.Println("n recommended events", len(recommendedEvents))
		if handleError(err) {
			return
		}
		setLoggedIn(recommendedEvents, loggedIn)
	}

	featuredEvents, err = db.GetTopEvents(15)
	if handleError(err) {
		return
	}
	setLoggedIn(featuredEvents, loggedIn)

	t, err := template.ParseFiles("../templates/index.html", "../templates/event_card.html")
	if handleError(err) {
		return
	}

	err = t.Execute(w, struct {
		LoggedIn          bool
		GoogleClientId    string
		FeaturedEvents    []EventCard
		RecommendedEvents []EventCard
	}{
		LoggedIn:          loggedIn,
		GoogleClientId:    os.Getenv("GOOGLE_CLIENT_ID"),
		FeaturedEvents:    featuredEvents,
		RecommendedEvents: recommendedEvents,
	})

	if handleError(err) {
		return
	}
}

func prefs(w http.ResponseWriter, req *http.Request) {

	loggedIn, email, id := loggedInEmail(req)
	if !loggedIn {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	fmt.Println(loggedIn, email, id)

	currentPrefs, err := db.GetPreferences(id)
	if err != nil {
		log.Println(err)
	}

	if req.Method == "GET" {
		respondWithTemplate(w, "../templates/prefs.html", currentPrefs)
	} else if req.Method == "PUT" {

		if err := req.ParseForm(); err != nil {
			log.Println(err)
		}

		for key, values := range req.PostForm {
			val := values[0]
			err := db.UpdatePreference(id, key, val)
			if err != nil {
				log.Println(err)
				fmt.Fprint(w, "Error")
				return
			}
			fmt.Println(key, val)
			fmt.Fprint(w, "Saved")
		}
	} else if req.Method == "DELETE" {
		fmt.Println("Deleting user")
		err := db.DeleteUser(id)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, "An error occurred. Please <a href=\"mailto:emilio@mywolverine.events\">email support</a>.")
		} else {
			http.SetCookie(w, &http.Cookie{Name: "jwt", Value: ""})
			fmt.Fprint(w, "Successfully deleted account. Close tab to exit.")
		}
	}

}

func vote(w http.ResponseWriter, req *http.Request) {

	loggedIn, _, id := loggedInEmail(req)
	params := req.URL.Query()
	voteType := params.Get("type")
	eventId, err := strconv.Atoi(params.Get("eventId"))

	if !loggedIn || err != nil {
		fmt.Fprint(w, "")
		return
	}
	fmt.Println(id, eventId, voteType)
	db.Vote(id, eventId, voteType)
	fmt.Fprint(w, "logged")
}

func main() {

	// Open the database connection
	var err error
	db, err = ConnectDB("../data/main.db")
	if err != nil {
		log.Fatal("Error connecting to db", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/prefs", prefs)
	mux.HandleFunc("/vote", vote)

	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("Static file server registered.")

	s := http.Server{
		Addr:         ":80",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	fmt.Println("Listening ...")
	err = s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
