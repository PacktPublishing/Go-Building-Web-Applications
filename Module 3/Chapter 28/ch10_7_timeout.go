package main

import (
	"fmt"
	"time"
)

func main() {

	myChan := make(chan int)

	go func() {
		time.Sleep(6 * time.Second)
		myChan <- 1
	}()

	for {
		select {
			case <-time.After(5 * time.Second):
				fmt.Println("This took too long!")
				return				
			case <-myChan:
				fmt.Println("Too little, too late")
		}
	}
}
