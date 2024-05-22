package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type EventCard struct {
	Title        string
	Subtitle     string
	Description  string
	EventID      string
	EventLink    string
	CalendarLink string
}

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

	featuredEvents := []EventCard{
		{"Test Event 1", "Subtitle Here", "Description here", "1", "https://google.com", "https://google.com"},
		{"Test Event 2", "Subtitle Here", "Description here", "2", "https://google.com", "https://google.com"},
		{"Test Event 3", "Subtitle Here", "Description here", "3", "https://google.com", "https://google.com"},
		{"Test Event 4", "Subtitle Here", "Description here", "4", "https://google.com", "https://google.com"},
	}

	t, err := template.ParseFiles("templates/index.html", "templates/event_card.html")

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
	}

	err = t.Execute(w, struct {
		LoggedIn       bool
		GoogleClientId string
		FeaturedEvents []EventCard
	}{
		LoggedIn:       loggedIn,
		GoogleClientId: os.Getenv("GOOGLE_CLIENT_ID"),
		FeaturedEvents: featuredEvents,
	})

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error")
	}
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

func vote(w http.ResponseWriter, req *http.Request) {

	loggedIn, email := loggedInEmail(req)
	if !loggedIn {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
	params := req.URL.Query()
	eventId := params.Get("eventId")
	responseStr := ""

	switch voteType := params.Get("type"); voteType {
	case "up":
		fmt.Println("up", eventId, email)
		responseStr = "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-hand-thumbs-up-fill\" viewBox=\"0 0 16 16\"><path d=\"M6.956 1.745C7.021.81 7.908.087 8.864.325l.261.066c.463.116.874.456 1.012.965.22.816.533 2.511.062 4.51a10 10 0 0 1 .443-.051c.713-.065 1.669-.072 2.516.21.518.173.994.681 1.2 1.273.184.532.16 1.162-.234 1.733q.086.18.138.363c.077.27.113.567.113.856s-.036.586-.113.856c-.039.135-.09.273-.16.404.169.387.107.819-.003 1.148a3.2 3.2 0 0 1-.488.901c.054.152.076.312.076.465 0 .305-.089.625-.253.912C13.1 15.522 12.437 16 11.5 16H8c-.605 0-1.07-.081-1.466-.218a4.8 4.8 0 0 1-.97-.484l-.048-.03c-.504-.307-.999-.609-2.068-.722C2.682 14.464 2 13.846 2 13V9c0-.85.685-1.432 1.357-1.615.849-.232 1.574-.787 2.132-1.41.56-.627.914-1.28 1.039-1.639.199-.575.356-1.539.428-2.59z\"/></svg>"
	case "down":
		fmt.Println("down", eventId, email)
		responseStr = "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-hand-thumbs-down-fill\" viewBox=\"0 0 16 16\"><path d=\"M6.956 14.534c.065.936.952 1.659 1.908 1.42l.261-.065a1.38 1.38 0 0 0 1.012-.965c.22-.816.533-2.512.062-4.51q.205.03.443.051c.713.065 1.669.071 2.516-.211.518-.173.994-.68 1.2-1.272a1.9 1.9 0 0 0-.234-1.734c.058-.118.103-.242.138-.362.077-.27.113-.568.113-.856 0-.29-.036-.586-.113-.857a2 2 0 0 0-.16-.403c.169-.387.107-.82-.003-1.149a3.2 3.2 0 0 0-.488-.9c.054-.153.076-.313.076-.465a1.86 1.86 0 0 0-.253-.912C13.1.757 12.437.28 11.5.28H8c-.605 0-1.07.08-1.466.217a4.8 4.8 0 0 0-.97.485l-.048.029c-.504.308-.999.61-2.068.723C2.682 1.815 2 2.434 2 3.279v4c0 .851.685 1.433 1.357 1.616.849.232 1.574.787 2.132 1.41.56.626.914 1.28 1.039 1.638.199.575.356 1.54.428 2.591\"/></svg>"
	case "cal":
		fmt.Println("cal", eventId, email)
		responseStr = "<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"16\" height=\"16\" fill=\"currentColor\" class=\"bi bi-calendar-check-fill\" viewBox=\"0 0 16 16\"><path d=\"M4 .5a.5.5 0 0 0-1 0V1H2a2 2 0 0 0-2 2v1h16V3a2 2 0 0 0-2-2h-1V.5a.5.5 0 0 0-1 0V1H4zM16 14V5H0v9a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2m-5.146-5.146-3 3a.5.5 0 0 1-.708 0l-1.5-1.5a.5.5 0 0 1 .708-.708L7.5 10.793l2.646-2.647a.5.5 0 0 1 .708.708\"/></svg>"
	default:
		responseStr = "what"
	}

	fmt.Fprint(w, responseStr)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/prefs", prefs)
	mux.HandleFunc("/vote", vote)

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
