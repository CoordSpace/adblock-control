package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

var apiKey *string
var piHoleURL *string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

// State is the contents of the newsAPI response data
type State struct {
	Status   string `json:"status"`
	Duration int    `json:"duration"`
}

func disableHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	params := u.Query()
	duration := params.Get("duration")
	endpoint := fmt.Sprintf("%s/api.php?disable=%s&auth=%s", *piHoleURL, duration, *apiKey)

	fmt.Println(endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	state := &State{}

	err = json.NewDecoder(resp.Body).Decode(&state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	minutes, err := strconv.Atoi(duration)

	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
		return
	}

	state.Duration = minutes / 60

	err = tpl.Execute(w, state)

	if err != nil {
		log.Println(err)
	}
}

func main() {
	// Get the URL, port, and API key from the command args if given
	apiKey = flag.String("apikey", "", "Pi-Hole API key.")
	piHoleURL = flag.String("url", "", "URL to the Pi-Hole Admin UI.")
	port := flag.String("port", "", "The port of the Adblock Control UI")
	flag.Parse()

	// Parse ENV variables if flags were not set.
	if *apiKey == "" {
		*apiKey = os.Getenv("API_KEY")
	}

	if *piHoleURL == "" {
		*piHoleURL = os.Getenv("URL")
	}

	if *port == "" {
		*port = os.Getenv("PORT")
		// Default to 8080 if nothing is set
		if *port == "" {
			*port = "8080"
		}
	}

	if *apiKey == "" || *piHoleURL == "" {
		log.Fatal("apiKey and url must be set via flag or ENV")
	}

	mux := http.NewServeMux()

	// import all our page assets
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/disable", disableHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe("0.0.0.0:"+(*port), mux)
}
