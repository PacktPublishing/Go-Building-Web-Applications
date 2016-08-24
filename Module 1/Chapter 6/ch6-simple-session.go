package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	var cookieStore = sessions.NewCookieStore([]byte("ideally, some random piece of entropy"))
	session, _ := cookieStore.Get(r, "mystore")
	if value, exists := session.Values["hello"]; exists {
		fmt.Fprintln(w, value)
	} else {
		session.Values["hello"] = "(world)"
		session.Save(r, w)
		fmt.Fprintln(w, "We just set the value!")
	}
}

func main() {
	http.HandleFunc("/test", cookieHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
