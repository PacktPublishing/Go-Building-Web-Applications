package main

import (
	"database/sql"
	"encoding/json"
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

type Comment struct {
	Id       int
	Name     string
	Email    string
	Comments string
}

type Page struct {
	Id         int
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Comments   []Comment
}

type JSONResponse struct {
	Fields map[string]string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	fmt.Println(pageGUID)
	err := database.QueryRow("SELECT id,page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Id, &thisPage.Title, &thisPage.RawContent, &thisPage.Date)
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
	APIOutput, _ := json.Marshal(thisPage)
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	fmt.Println(APIOutput)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, thisPage)
}

func APICommentPost(w http.ResponseWriter, r *http.Request) {
	var commentAdded string
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error)
	}
	fmt.Println(r.FormValue)
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")

	res, err := database.Exec("INSERT INTO comments SET comment_name=?, comment_email=?, comment_text=?", name, email, comments)

	if err != nil {
		log.Println(err.Error)
	}

	id, err := res.LastInsertId()
	if err != nil {
		commentAdded = "false"
	} else {
		commentAdded = "true"
	}

	var resp JSONResponse
	resp.Fields["id"] = string(id)
	resp.Fields["added"] = commentAdded
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonResp)
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
	routes.HandleFunc("/api/comments", APICommentPost).
		Methods("POST")
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	http.Handle("/", routes)
	http.ListenAndServe(PORT, nil)

}
