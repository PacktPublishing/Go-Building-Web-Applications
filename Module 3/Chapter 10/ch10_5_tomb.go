package main

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/tomb"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var URLS []URL

type GoTomb struct {
	tomb tomb.Tomb
	wg   sync.WaitGroup
}

type URL struct {
	Status bool
	URI    string
	Body   string
}

func (gt GoTomb) Kill() {

	gt.tomb.Kill(nil)

}

func (gt *GoTomb) TombListen(ii int) {

	for {
		select {
		case <-gt.tomb.Dying():
			fmt.Println("Got kill command from tomb!")
			if URLS[ii].Status == false {
				fmt.Println("Never got data for", URLS[ii].URI)
			}
			return
		}
	}
}

func (gt *GoTomb) Fetch() {
	for i := range URLS {
		go gt.TombListen(i)

		go func(ii int) {

			timeDelay := 5 * ii
			fmt.Println("Waiting ", strconv.FormatInt(int64(timeDelay), 10), " seconds to get", URLS[ii].URI)
			time.Sleep(time.Duration(timeDelay) * time.Second)
			response, _ := http.Get(URLS[ii].URI)
			URLS[ii].Status = true
			fmt.Println("Got body for ", URLS[ii].URI)
			responseBody, _ := ioutil.ReadAll(response.Body)
			URLS[ii].Body = string(responseBody)
		}(i)
	}
}

func main() {

	done := make(chan int)

	URLS = []URL{{Status: false, URI: "http://www.google.com", Body: ""}, {Status: false, URI: "http://www.amazon.com", Body: ""}, {Status: false, URI: "http://www.ubuntu.com", Body: ""}}

	var MasterChannel GoTomb
	MasterChannel.Fetch()

	go func() {

		time.Sleep(10 * time.Second)
		MasterChannel.Kill()
		done <- 1
	}()

	for {
		select {
		case <-done:
			fmt.Println("")
			return
		default:
		}
	}
}
