package main

import
(
	"sync"
	"fmt"
)

type Fetcher struct {
	Done bool
	InChannel chan string
	OutChannel chan string
	Parser Parser
}

func (f Fetcher) Get(url string, fetchParser Parser) {
	f.Done = true
	f.Parser = fetchParser
	f.OutChannel <- url
}

type Parser struct {
	Done bool
	InChannel chan string
	OutChannel chan Item
}

func (p Parser) Listen () {
	for {
		select {


		}
	}
}

type Item struct {

}

type URL struct {
	Mutex sync.Mutex
	URI string
}

type Feed struct {
	Mutex sync.Mutex
	Name string
	URI string
}

var Feeds []Feed

func main() {

	Feeds := []Feed{ { Name: "Google News", URI: "http://news.google.com/?output=rss"}, { Name: "Golang Nuts", URI: "https://groups.google.com/forum/feed/Golang-nuts/msgs/rss_v2_0.xml?num=50"} }

	for i := range Feeds {
		fmt.Println(i,"...")
	}
}