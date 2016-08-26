package main

import
(
"fmt"

)

type Item struct {
	Url string
	Data []byte
}

type Feed struct {
	Url string
	Name string
	Items []Item
}

var Feeds []Feed


func process(cM master) {

	for _,i := range Feeds {

		fmt.Println("feed",i)
		item := Item{}
		item.Url = i.Url
		cM <- item
	}

}

func processItem(url string) {
	
}

type master chan Item

func main() {

	done := make(chan bool)

	Feeds = []Feed{ Feed{ Name: "New York Times", Url: "http://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml"}, Feed{ Name: "Wall Street Journal", Url: "http://feeds.wsjonline.com/wsj/xml/rss/3_7011.xml"} }
	feedChannel := make(master)
	

	go process(feedChannel)

	select {
		case fm := <-feedChannel:
			fmt.Println("Got URL",fm.Url)
			processItem(fm.Url)
	}

	<- done
	fmt.Println("Done!")

}