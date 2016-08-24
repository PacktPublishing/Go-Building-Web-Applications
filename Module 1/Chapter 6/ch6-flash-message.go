package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}

var (
	templates = template.Must(template.ParseGlob("templates/*"))
	port      = ":8080"
)

func startHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "ch6-flash.html", nil)
	if err != nil {
		log.Fatal("Template ch6-flash missing")
	}
}

func middleHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue := r.PostFormValue("message")
	cookie := http.Cookie{Name: "message", Value: "message:" + cookieValue, Expires: time.Now().Add(60 * time.Second), HttpOnly: true}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/finish", 301)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	cookieVal, _ := r.Cookie("message")

	if cookieVal != nil {
		fmt.Fprintln(w, "We found: "+string(cookieVal.Value)+", but try to refresh!")
		cookie := http.Cookie{Name: "message", Value: "", Expires: time.Now(), HttpOnly: true}
		http.SetCookie(w, &cookie)
	} else {
		fmt.Fprintln(w, "That cookie was gone in a flash")
	}

}
func main() {

	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/middle", middleHandler)
	http.HandleFunc("/finish", finishHandler)
	log.Fatal(http.ListenAndServe(port, nil))

}
