package main

import (
	"log"
	"time"
)

func main() {

	timeout := time.NewTimer(5 * time.Second)
	defer log.Println("Timed out!")

	for {
		select {
		case <-timeout.C:
			return
		default:
		}
	}

}
