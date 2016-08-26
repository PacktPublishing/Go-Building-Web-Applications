package main

import (
	"fmt"
)

func main() {

	var channel chan int

	if channel != nil {
		channel <- 1	
	}

	for {
		select {
			case <- channel: 
				
			default:
		}
	}

}
