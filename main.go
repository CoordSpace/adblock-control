package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

var app_pass *string
var piHoleURL *string

type Login struct {
	Password string `json:"password"`
}

type DNS struct {
	Blocking bool   `json:"blocking"`
	Timer    int    `json:"timer"`
	SID      string `json:"sid"`
}

type DNSReply struct {
	Blocking string  `json:"blocking"`
	Timer    int     `json:"timer"`
	Took     float32 `json:"took"`
}

type Session struct {
	Valid    bool   `json:"valid"`
	TOTP     bool   `json:"totp"`
	SID      string `json:"sid"`
	CSRF     string `json:"csrf"`
	Validity int    `json:"validity"`
	Message  string `json:"message"`
}

type SessionWrapper struct {
	Session Session `json:"session"`
	Took    float32 `json:"took"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func disableHandler(w http.ResponseWriter, r *http.Request) {

	u, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	// parse duration param from the site form
	params := u.Query()

	// Create the client we'll be using for all the pihole connections
	// pihole default config generates a self-signed cert, so disable verify
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Get the SID
	login := &Login{Password: *app_pass}

	message_data, err := json.Marshal(login)

	if err != nil {
		println("Error marshalling json data!")
		log.Fatal(err)
	}

	authRequestURL := fmt.Sprintf("%s/auth", *piHoleURL)

	req, err := http.NewRequest("POST", authRequestURL, bytes.NewBuffer(message_data))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		println("Error creating POST request!")
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		println("Error sending auth POST to pihole!")
		log.Fatal(err)
	}

	var s SessionWrapper
	// convert from http response reader to bytes
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		println("Error parsing auth response from pihole!")
		log.Fatal(err)
	}

	// sid is now in session.SID
	err = json.Unmarshal(body, &s)

	if err != nil {
		println("Error converting auth response from pihole to JSON!")
		log.Fatal(err)
	}

	if !s.Session.Valid {
		log.Fatal("Error getting valid session from API!")
	}

	duration, err := strconv.Atoi(params.Get("duration"))

	if err != nil {
		println("Error converting duration to int!")
		log.Fatal(err)
	}

	// Disable the adblock
	dnsRequestURL := fmt.Sprintf("%s/dns/blocking", *piHoleURL)

	dns := DNS{Blocking: false, Timer: duration, SID: s.Session.SID}

	dns_data, err := json.Marshal(dns)

	if err != nil {
		println("Error marshalling json data!")
		log.Fatal(err)
	}

	req, err = http.NewRequest("POST", dnsRequestURL, bytes.NewBuffer(dns_data))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		println("Error creating POST request!")
		log.Fatal(err)
	}

	resp, err = client.Do(req)
	if err != nil {
		println("Error sending dns POST to pihole!")
		log.Fatal(err)
	}

	var dns_result DNSReply

	body, err = io.ReadAll(resp.Body)

	if err != nil {
		println("Error reading response body!")
		log.Fatal(err)
	}

	// convert the reply from pihole into a DNS status struct
	err = json.Unmarshal(body, &dns_result)

	if err != nil {
		println("Error parsing DNS JSON from Pi-hole!")
		log.Fatal(err)
	}

	// convert to minutes for easier reading on the result page
	dns_result.Timer = dns_result.Timer / 60

	// print out information from the ack body
	err = tpl.Execute(w, dns_result)

	if err != nil {
		println("Error executing the DNS result to the page!")
		log.Fatal(err)
	}
}

func main() {
	// Get the URL, port, and API key from the command args, if given.
	app_pass = flag.String("app_pass", "", "Pi-Hole API app password.")
	piHoleURL = flag.String("url", "", "URL to the Pi-Hole Admin UI.")
	port := flag.String("port", "", "The port of the Adblock Control UI")
	flag.Parse()

	// Parse ENV variables if flags were not set.
	if *app_pass == "" {
		*app_pass = os.Getenv("APP_PASS")
	}

	if *piHoleURL == "" {
		*piHoleURL = os.Getenv("URL")
	}

	if *port == "" {
		*port = os.Getenv("PORT")
		// Default the app's port to 8080 if nothing is set
		if *port == "" {
			*port = "8080"
		}
	}

	if *app_pass == "" || *piHoleURL == "" {
		log.Fatal("app_pass and url must be set via flag or ENV")
	}

	mux := http.NewServeMux()

	// import all our page assets
	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/disable", disableHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe("0.0.0.0:"+(*port), mux)
}
