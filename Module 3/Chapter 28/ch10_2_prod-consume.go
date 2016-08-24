package main

import (
	"fmt"
	"log"
	"time"
)

const CONSUMERS = 5

func main() {

	Producer := make(chan (chan int))

	for i := 0; i < CONSUMERS; i++ {
		go func() {
			time.Sleep(1000 * time.Microsecond)
			conChan := make(chan int)

			go func() {
				for {
					select {
					case _, ok := <-conChan:
						if ok == false {
							return
						} else {
							Producer <- conChan
						}
					default:
					}
				}
			}()

			conChan <- 1
			close(conChan)
		}()
	}
	for {
		select {
		case consumer, ok := <-Producer:
			if ok == false {
				fmt.Println("Goroutine closed?")
				close(Producer)
			} else {
				log.Println(consumer)
				// consumer <- 1
			}
			fmt.Println("Got message from secondary channel")
		default:
		}
	}
}
