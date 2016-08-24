package main

import (
	"net/http"
)

const (
	PORT = ":8080"
)

func main() {

	http.ListenAndServe(PORT, http.FileServer(http.Dir("/var/www")))
}
