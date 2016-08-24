package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

func HTMLHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	fmt.Fprintln(w, input)
}

func HTMLHandlerSafe(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	input = template.HTMLEscapeString(input)
	fmt.Fprintln(w, input)
}

func JSHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	fmt.Fprintln(w, input)
}

func JSHandlerSafe(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	input = template.JSEscapeString(input)
	fmt.Fprintln(w, input)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/html", HTMLHandler)
	router.HandleFunc("/js", JSHandler)
	router.HandleFunc("/html_safe", HTMLHandlerSafe)
	router.HandleFunc("/js_safe", JSHandlerSafe)
	http.ListenAndServe(":8080", router)
}
