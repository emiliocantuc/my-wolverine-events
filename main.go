package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func respondWithTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		t.Execute(w, data)
	}
}

func index(w http.ResponseWriter, req *http.Request) {

	loggedIn, email := loggedInEmail(req)
	fmt.Println(loggedIn, email)

	// Checked if logged in and respond with template
	respondWithTemplate(w, "templates/index.html", struct {
		LoggedIn       bool
		GoogleClientId string
	}{
		LoggedIn:       loggedIn,
		GoogleClientId: os.Getenv("GOOGLE_CLIENT_ID"),
	})

}

func prefs(w http.ResponseWriter, req *http.Request) {

	loggedIn, email := loggedInEmail(req)
	if !loggedIn {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	fmt.Println(loggedIn, email)

	templateStruct := struct {
		EducationLevel  string
		Program         string
		CampusLocation  string
		Interests       string
		IncSeminars     bool
		IncSports       bool
		IncSocial       bool
		KeywordsToAvoid string
	}{
		EducationLevel: "Graduate", Program: "Statistics", CampusLocation: "Central Campus",
		Interests: "machine learning, big data, soccer", IncSeminars: true, KeywordsToAvoid: "housing, job"}

	if req.Method == "GET" {
		respondWithTemplate(w, "templates/prefs.html", templateStruct)
	} else if req.Method == "PUT" {

		if err := req.ParseForm(); err != nil {
			log.Println(err)
		}

		for key, values := range req.PostForm {
			val := values[0]
			// update in db here
			fmt.Println(key, val)
			fmt.Fprint(w, "Saved")
		}
	}

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/prefs", prefs)

	s := http.Server{
		Addr:         ":80",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	fmt.Println("Listening ...")
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
