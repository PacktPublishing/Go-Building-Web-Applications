package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/streadway/amqp"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	DBHost  = "127.0.0.1"
	DBPort  = ":3306"
	DBUser  = "root"
	DBPass  = ""
	DBDbase = "cms"
	PORT    = ":8080"
	MQHost  = "127.0.0.1"
	MQPort  = ":5672"
)

var database *sql.DB
var sessionStore = sessions.NewCookieStore([]byte("our-social-network-application"))
var UserSession Session

var WelcomeTitle = "You've succcessfully registered!"
var WelcomeEmail = "Welcome to our CMS, {{Email}}!  We're glad you could join us."

type RegistrationData struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type Comment struct {
	Id          int
	Name        string
	Email       string
	CommentText string
}

type Page struct {
	Id         int
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Comments   []Comment
	Session    Session
	GUID       string
}

type User struct {
	Id   int
	Name string
}

type Session struct {
	Id              string
	Authenticated   bool
	Unauthenticated bool
	User            User
}

type JSONResponse struct {
	Fields map[string]string
}

func getSessionUID(sid string) int {
	user := User{}
	err := database.QueryRow("SELECT user_id FROM sessions WHERE session_id=?", sid).Scan(user.Id)
	if err != nil {
		fmt.Println(err.Error)
		return 0
	}
	return user.Id
}

func updateSession(sid string, uid int) {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	_, err := database.Exec("INSERT INTO sessions SET session_id=?, user_id=?, session_update=? ON DUPLICATE KEY UPDATE user_id=?, session_update=?", sid, uid, tstamp, uid, tstamp)
	if err != nil {
		fmt.Println(err.Error)
	}
}

func generateSessionId() string {
	sid := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, sid)
	if err != nil {
		log.Fatal("Could not generate session id")
	}
	return base64.URLEncoding.EncodeToString(sid)
}

func validateSession(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, "app-session")
	if sid, valid := session.Values["sid"]; valid {
		currentUID := getSessionUID(sid.(string))
		updateSession(sid.(string), currentUID)
		UserSession.Id = string(currentUID)
	} else {
		newSID := generateSessionId()
		session.Values["sid"] = newSID
		session.Save(r, w)
		UserSession.Id = newSID
		updateSession(newSID, 0)
	}
	fmt.Println(session.ID)
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	validateSession(w, r)
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	thisPage.GUID = pageGUID
	err := database.QueryRow("SELECT id,page_title,page_content,page_date FROM pages WHERE page_guid=?", pageGUID).Scan(&thisPage.Id, &thisPage.Title, &thisPage.RawContent, &thisPage.Date)
	thisPage.Content = template.HTML(thisPage.RawContent)
	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(err)
		return
	}

	comments, err := database.Query("SELECT id, comment_name as Name, comment_email, comment_text FROM comments WHERE page_id=?", thisPage.Id)
	if err != nil {
		log.Println(err)
	}
	for comments.Next() {
		var comment Comment
		comments.Scan(&comment.Id, &comment.Name, &comment.Email, &comment.CommentText)
		thisPage.Comments = append(thisPage.Comments, comment)
	}
	thisPage.Session.Authenticated = false
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

func APICommentPut(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Heyyyyyyy")
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error)
	}
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(id)
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")
	fmt.Println("UPDATE comments SET comment_name=?, comment_email=?, comment_text=? WHERE comment_id=?", name, email, comments, id)
	res, err := database.Exec("UPDATE comments SET comment_name=?, comment_email=?, comment_text=? WHERE comment_id=?", name, email, comments, id)
	fmt.Println(res)
	if err != nil {
		log.Println(err.Error)
	}

	var resp JSONResponse

	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, jsonResp)
}

func weakPasswordHash(password string) []byte {
	hash := sha1.New()
	io.WriteString(hash, password)
	return hash.Sum(nil)
}

func MQConnect() (*amqp.Connection, *amqp.Channel, error) {
	url := "amqp://" + MQHost + MQPort
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	if _, err := channel.QueueDeclare("", false, true, false, false, nil); err != nil {
		return nil, nil, err
	}
	return conn, channel, nil
}

func MQPublish(message []byte) {
	err = channel.Publish(
		"email", // exchange
		"",      // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err.Error)
	}
	name := r.FormValue("user_name")
	email := r.FormValue("user_email")
	pass := r.FormValue("user_password")
	pageGUID := r.FormValue("referrer")
	// pass2 := r.FormValue("user_password2")
	gure := regexp.MustCompile("[^A-Za-z0-9]+")
	guid := gure.ReplaceAllString(name, "")
	password := weakPasswordHash(pass)

	res, err := database.Exec("INSERT INTO users SET user_name=?, user_guid=?, user_email=?, user_password=?", name, guid, email, password)
	fmt.Println(res)
	if err != nil {
		fmt.Fprintln(w, err.Error)
	} else {
		Email := RegistrationData{Email: email, Message: ""}
		message, err := template.New("email").Parse(WelcomeEmail)
		var mbuf bytes.Buffer
		message.Execute(&mbuf, Email)
		MQPublish(json.Marshal(mbuf.String()))
		http.Redirect(w, r, "/page/"+pageGUID, 301)
	}
}
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	validateSession(w, r)
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err.Error)
	}
	u := User{}
	name := r.FormValue("user_name")
	pass := r.FormValue("user_password")
	password := weakPasswordHash(pass)
	err = database.QueryRow("SELECT user_id, user_name FROM users WHERE user_name=? and user_password=?", name, password).Scan(&u.Id, &u.Name)
	if err != nil {
		fmt.Fprintln(w, err.Error)
		u.Id = 0
		u.Name = ""
	} else {
		updateSession(UserSession.Id, u.Id)
		fmt.Fprintln(w, u.Name)
	}
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
	routes.HandleFunc("/register", RegisterPOST).Methods("POST")
	routes.HandleFunc("/login", LoginPOST).
		Methods("POST")
	routes.HandleFunc("/api/pages", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/page/{id:[\\w\\d\\-]+}", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/comments", APICommentPost).
		Methods("POST")
	routes.HandleFunc("/api/comments/{id:[\\w\\d\\-]+}", APICommentPut).
		Methods("PUT")
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
	http.Handle("/", routes)
	http.ListenAndServe(PORT, nil)

}
