package main

import (
	"fmt"
	"net/http"
	"time"
)

const INIT_DELAY = 3000
const MAX_DELAY = 60000
const MAX_RETRIES = 4
const DELAY_INCREMENT = 5000

var Servers []Server

type Server struct {
	Name        string
	URI         string
	LastChecked time.Time
	Status      bool
	StatusCode  int
	Delay       int
	Retries     int
	Channel     chan bool
}

func (s *Server) checkServerStatus(sc chan *Server) {
	var previousStatus string

	if s.Status == true {
		previousStatus = "OK"
	} else {
		previousStatus = "down"
	}

	fmt.Println("Checking Server", s.Name)
	fmt.Println("\tServer was", previousStatus, "on last check at", s.LastChecked)

	response, err := http.Get(s.URI)
	if err != nil {
		fmt.Println("\tError: ", err)
		s.Status = false
		s.StatusCode = 0
	} else {
		fmt.Println(response.Status)
		s.StatusCode = response.StatusCode
		s.Status = true
	}

	s.LastChecked = time.Now()
	sc <- s
}

func cycleServers(sc chan *Server) {

	for i := 0; i < len(Servers); i++ {
		Servers[i].Channel = make(chan bool)
		go Servers[i].updateDelay(sc)
		Servers[i].checkServerStatus(sc)
	}
}

func (s *Server) updateDelay(sc chan *Server) {
	for {
		select {
		case msg := <-s.Channel:

			if msg == false {
				s.Delay = s.Delay + DELAY_INCREMENT
				s.Retries++
				if s.Delay > MAX_DELAY {
					s.Delay = MAX_DELAY
				}

			} else {
				s.Delay = INIT_DELAY
			}
			newDuration := time.Duration(s.Delay)

			if s.Retries <= MAX_RETRIES {
				fmt.Println("\tWill check server again")
				time.Sleep(newDuration * time.Millisecond)
				s.checkServerStatus(sc)
			} else {
				fmt.Println("\tServer not reachable after", MAX_RETRIES, "retries")
			}

		default:
		}
	}
}

func main() {
	fmt.Println("")
	fmt.Println("")	
	endChan := make(chan bool)
	serverChan := make(chan *Server)

	Servers = []Server{{Name: "Google", URI: "http://www.google.com", Status: true, Delay: INIT_DELAY}, {Name: "Yahoo", URI: "http://www.yahoo.com", Status: true, Delay: INIT_DELAY}, {Name: "Amazon", URI: "http://amazon.zom", Status: true, Delay: INIT_DELAY}}

	go cycleServers(serverChan)

	for {
		select {
		case currentServer := <-serverChan:
			currentServer.Channel <- false
		default:

		}
	}

	<-endChan

}
