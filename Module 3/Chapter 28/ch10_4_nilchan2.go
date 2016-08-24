package main

import
(
	"log"
	"time"
)

func main() {
	
	done := make(chan int)
	defer close(done)
	defer log.Println("End of script")
	go func() {
		time.Sleep(time.Second * 5)
		done <- 1
	}()

	for {
		select {
			case <- done: 
				log.Println("Got transmission")
				return
			default:
		}
	}


}