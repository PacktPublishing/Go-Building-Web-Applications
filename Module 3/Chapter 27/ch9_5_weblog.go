package main

import (
	"code.google.com/p/log4go"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

var errorLog log4go.Logger
var errorLogWriter log4go.FileLogWriter

var accessLog log4go.Logger
var accessLogWriter *log4go.FileLogWriter

var screenLog log4go.Logger

var networkLog log4go.Logger

func init() {
	fmt.Println("Web Server Starting")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	pageFoundMessage := "Page found: " + r.URL.Path
	accessLog.Info(pageFoundMessage)
	networkLog.Info(pageFoundMessage)
 	w.Write([]byte("Valid page"))		
}

func notFound(w http.ResponseWriter, r *http.Request) {
	pageNotFoundMessage := "Page not found / 404: " + r.URL.Path
	errorLog.Info(pageNotFoundMessage)
 	w.Write([]byte("Page not found"))	
}

func restricted(w http.ResponseWriter, r *http.Request) {
	message := "Restricted directory access attempt!"
	errorLog.Info(message)
	accessLog.Info(message)
	screenLog.Info(message)
	networkLog.Info(message)
	w.Write([]byte("Restricted!"))

}

func main() {

	screenLog = make(log4go.Logger)
	screenLog.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

 	errorLogWriter := log4go.NewFileLogWriter("web-errors.log", false)
    errorLogWriter.SetFormat("%d %t - %M (%S)")
    errorLogWriter.SetRotate(false)
    errorLogWriter.SetRotateSize(0)
    errorLogWriter.SetRotateLines(0)
    errorLogWriter.SetRotateDaily(true)

	errorLog = make(log4go.Logger)
	errorLog.AddFilter("file", log4go.DEBUG, errorLogWriter)

	networkLog = make(log4go.Logger)
	networkLog.AddFilter("network", log4go.DEBUG, log4go.NewSocketLogWriter("tcp", "localhost:3000"))

	accessLogWriter = log4go.NewFileLogWriter("web-access.log",false)
    accessLogWriter.SetFormat("%d %t - %M (%S)")
    accessLogWriter.SetRotate(true)
    accessLogWriter.SetRotateSize(0)
    accessLogWriter.SetRotateLines(500)
    accessLogWriter.SetRotateDaily(false)	

   	accessLog = make(log4go.Logger)
   	accessLog.AddFilter("file",log4go.DEBUG,accessLogWriter) 

	rtr := mux.NewRouter()
	rtr.HandleFunc("/valid", pageHandler)
	rtr.HandleFunc("/.git/", restricted)	
	rtr.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", rtr)
	http.ListenAndServe(":8080", nil)
}
