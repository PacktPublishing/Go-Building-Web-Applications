package main

import
(
"net/http"
"html/template"
"time"
"regexp"
"fmt"
"io/ioutil"
"database/sql"
"log"
"runtime"
_ "github.com/go-sql-driver/mysql"
)

const staticPath string = "static/"

type WebPage struct {

	Title string
	Contents string
	Connection *sql.DB

}

type customRouter struct {

}

func serveDynamic() {

}

func serveRendered() {

}

func serveStatic() {

}

func (customRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path;

	staticPatternString := "static/(.*)"
	templatePatternString := "template/(.*)"
	dynamicPatternString := "dynamic/(.*)"

	staticPattern := regexp.MustCompile(staticPatternString)
	templatePattern := regexp.MustCompile(templatePatternString)
	dynamicDBPattern := regexp.MustCompile(dynamicPatternString)

	if staticPattern.MatchString(path) {
		 serveStatic()
		page := staticPath + staticPattern.ReplaceAllString(path, "${1}") + ".html"
		http.ServeFile(rw, r, page)
	}else if templatePattern.MatchString(path) {
		
		serveRendered()
		urlVar := templatePattern.ReplaceAllString(path, "${1}")

		page.Title = "This is our URL: " + urlVar
		customTemplate.Execute(rw,page)
		
	}else if dynamicDBPattern.MatchString(path) {
		
		serveDynamic()
		page = getArticle(1)
		customTemplate.Execute(rw,page)
	}

}

func gobble(s []byte) {


}


var customHTML string
var customTemplate template.Template
var page WebPage
var templateSet bool
var Database sql.DB

func getArticle(id int) WebPage {
	Database,err := sql.Open("mysql", "test:test@/master")				
	if err != nil {
	   fmt.Println("DB error!")
	}		

	var articleTitle string
	sqlQ := Database.QueryRow("SELECT article_title from articles where article_id=? LIMIT 1", id).Scan(&articleTitle)
	switch {
		case sqlQ == sql.ErrNoRows:
		    fmt.Printf("No rows!")
		case sqlQ != nil:
		    fmt.Println(sqlQ)
		default:
		  
	}	


	wp := WebPage{}
	wp.Title = articleTitle
	return wp

}

func main() {

	runtime.GOMAXPROCS(2)

	var cr customRouter;

	fileName := staticPath + "template.html"
	cH,_ := ioutil.ReadFile(fileName)
	customHTML = string(cH[:])

	page := WebPage{ Title: "This is our URL: ", Contents: "Enjoy our content" }
	cT,_ := template.New("Hey").Parse(customHTML)
	customTemplate = *cT

	gobble(cH)
	log.Println(page)
	fmt.Println(customTemplate)


	server := &http.Server {
			Addr: ":9000",
			Handler:cr,
			ReadTimeout: 10 * time.Second,
			WriteTimeout: 10 * time.Second,
			MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()

}