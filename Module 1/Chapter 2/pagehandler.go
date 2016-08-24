package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

const (
	PORT = ":8080"
)

func pageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	fileName := "files/" + pageID + ".html"
	_, err := os.Stat(fileName)
	if err != nil {
		fileName = "files/404.html"
	}

	http.ServeFile(w, r, fileName)
}
