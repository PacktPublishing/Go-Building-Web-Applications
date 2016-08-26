package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

const MAX_SERVER_FAILURES = 10
const DEFAULT_TIMEOUT_SECONDS = 5
const MAX_TIMEOUT_SECONDS = 60
const TIMEOUT_INCREMENT = 5
const MAX_RETRIES = 5

type Server struct {
	Name        string
	Failures    int
	InService   bool
	Status      bool
	StatusCode  int
	Addr        string
	Timeout     int
	LastChecked time.Time
	Recheck     chan bool
}

func (s *Server) serverListen(sc chan bool) {
	for {
		select {
		case msg := <-s.Recheck:
			var statusText string
			if msg == false {
				statusText = "NOT in service"
				s.Failures++
				s.Timeout = s.Timeout + TIMEOUT_INCREMENT
				if s.Timeout > MAX_TIMEOUT_SECONDS {
					s.Timeout = MAX_TIMEOUT_SECONDS
				}
			} else {
				if ServersAvailable == false {
					ServersAvailable = true
					sc <- true
				}
				statusText = "in service"
				s.Timeout = DEFAULT_TIMEOUT_SECONDS
			}

			if s.Failures >= MAX_SERVER_FAILURES {
				s.InService = false
				fmt.Println("\tServer", s.Name, "failed too many times.")
			} else {
				timeString := strconv.FormatInt(int64(s.Timeout), 10)
				fmt.Println("\tServer", s.Name, statusText, "will check again in", timeString, "seconds")
				s.InService = true
				time.Sleep(time.Second * time.Duration(s.Timeout))
				go s.checkStatus()
			}

		}
	}
}

func (s *Server) checkStatus() {
	previousStatus := "Unknown"
	if s.Status == true {
		previousStatus = "OK"
	} else {
		previousStatus = "down"
	}
	fmt.Println("Checking Server", s.Name)
	fmt.Println("\tServer was", previousStatus, "on last check at", s.LastChecked)
	response, err := http.Get(s.Addr)
	if err != nil {
		fmt.Println("\tError: ", err)
		s.Status = false
		s.StatusCode = 0
	} else {
		s.StatusCode = response.StatusCode
		s.Status = true
	}

	s.LastChecked = time.Now()
	s.Recheck <- s.Status
}

func healthCheck(sc chan bool) {
	fmt.Println("Running initial health check")
	for i := range Servers {
		Servers[i].Recheck = make(chan bool)
		go Servers[i].serverListen(sc)
		go Servers[i].checkStatus()
	}
}

func roundRobin() Server {
	var AvailableServer Server

	if nextServerIndex > (len(Servers) - 1) {
		nextServerIndex = 0
	}

	if Servers[nextServerIndex].InService == true {
		AvailableServer = Servers[nextServerIndex]
	} else {
		serverReady := false
		for serverReady == false {
			for i := range Servers {
				if Servers[i].InService == true {
					AvailableServer = Servers[i]
					serverReady = true
				}
			}

		}
	}
	nextServerIndex++
	return AvailableServer
}

var Servers []Server
var nextServerIndex int
var ServersAvailable bool
var ServerChan chan bool
var Proxy *httputil.ReverseProxy
var ResetProxy chan bool

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	Proxy = setProxy()
	return func(w http.ResponseWriter, r *http.Request) {

		r.URL.Path = "/"

		p.ServeHTTP(w, r)

	}
}

func setProxy() *httputil.ReverseProxy {

	nextServer := roundRobin()
	nextURL, _ := url.Parse(nextServer.Addr)
	log.Println("Next proxy source:", nextServer.Addr)
	prox := httputil.NewSingleHostReverseProxy(nextURL)

	return prox
}

func startListening() {
	http.HandleFunc("/index.html", handler(Proxy))
	_ = http.ListenAndServe(":8080", nil)

}

func main() {
	nextServerIndex = 0
	ServersAvailable = false
	ServerChan := make(chan bool)
	done := make(chan bool)

	fmt.Println("Starting load balancer")
	Servers = []Server{{Name: "Web Server 01", Addr: "http://www.google.com", Status: false, InService: false}, {Name: "Web Server 02", Addr: "http://www.amazon.com", Status: false, InService: false}, {Name: "Web Server 03", Addr: "http://www.apple.zom", Status: false, InService: false}}

	go healthCheck(ServerChan)

	for {
		select {
		case <-ServerChan:
			Proxy = setProxy()
			startListening()
			return

		}
	}

	/*

	*/

	<-done
}
