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

func handleInternalServerError(w http.ResponseWriter, err error, msg string) {
	if msg == "" {
		msg = "Error interno del servidor"
	}
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

func respondWithTemplate(w http.ResponseWriter, data interface{}, tmpls ...string) {
	t, err := template.ParseFiles(tmpls...)
	if err != nil {
		handleInternalServerError(w, err, "")
		return
	}
	if err = t.Execute(w, data); err != nil {
		handleInternalServerError(w, err, "")
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
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: jwtString, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode, MaxAge: 30 * 24 * 60 * 60})
	http.Redirect(w, req, redirectTo, http.StatusSeeOther)
}

func logout(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: ""})
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, req *http.Request) {

	var featuredEvents, recommendedEvents []EventCard
	var err error

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
		if err != nil {
			handleInternalServerError(w, err, "")
			return
		}
		setLoggedIn(recommendedEvents, loggedIn)
	}

	featuredEvents, err = db.GetTopEvents(15)
	fmt.Println("n top events", len(featuredEvents))
	if err != nil {
		handleInternalServerError(w, err, "")
		return
	}
	setLoggedIn(featuredEvents, loggedIn)

	respondWithTemplate(w, struct {
		LoggedIn          bool
		GoogleClientId    string
		FeaturedEvents    []EventCard
		RecommendedEvents []EventCard
	}{
		LoggedIn:          loggedIn,
		GoogleClientId:    os.Getenv("GOOGLE_CLIENT_ID"),
		FeaturedEvents:    featuredEvents,
		RecommendedEvents: recommendedEvents,
	}, "../templates/index.html", "../templates/event_card.html")
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
		respondWithTemplate(w, currentPrefs, "../templates/prefs.html")
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
	db, err = ConnectDB("../data/main.db?_cache_size=0")
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
	fs := http.FileServer(http.Dir("../static"))
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
