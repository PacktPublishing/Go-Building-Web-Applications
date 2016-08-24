package main

import (
	"crypto/tls"
	"database/sql"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = ""
	DBDbase = "cms"
	PORT    = ":8080"
)

var database *sql.DB

type Page struct {
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err)
		return
	}
	// html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`

	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func APIPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, thisPage)
}

func main() {
	dbConn := fmt.Sprintf("%s:%s@/%s", DBUser, DBPass, DBDbase)
	fmt.Println(dbConn)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("Couldn't connect!")
		log.Println(err.Error)
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/api/pages", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/page/{id:[\\w\\d\\-]+}", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/register", RegisterPOST).
		Methods("POST").
		Schemes("https")
	routes.HandleFunc("/login", LoginPOST).
		Methods("POST").
		Schemes("https")
	http.Handle("/", routes)
	certificates, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	tlsConf := tls.Config{Certificates: []tls.Certificate{certificates}}
	tls.Listen("tcp", PORT, &tlsConf)

}
