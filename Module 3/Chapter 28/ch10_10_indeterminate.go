package main

import (

	"fmt"
	"time"
)

func main() {

	acceptingChannel := make(chan interface{})

	go func() {

		acceptingChannel <- "A text message"
		time.Sleep(3 * time.Second)
		acceptingChannel <- false
	}()

	for {
		select {
			case msg := <- acceptingChannel:
				switch typ := msg.(type) {
					case string:
						fmt.Println("Got text message",typ)
					case bool: 
						fmt.Println("Got boolean message",typ)
						if typ == false {
							return
						}
					default:
						fmt.Println("Some other type of message")
				}
				
			default:

		}

	}

	 <- acceptingChannel
}
