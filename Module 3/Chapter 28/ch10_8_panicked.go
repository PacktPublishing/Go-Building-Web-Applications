package main

import
(
"fmt"
)

const TOTAL_RANDOMS = 100

func sendRandoms(ch chan int) {
	for i := 0; i < TOTAL_RANDOMS; i++ {
		ch <- i
	}
}

func main() {

	ch := make(chan int)

	go sendRandoms(ch)

	for {
		select {
			case num := <- ch:
				fmt.Println(num)
				if num == 98 {
					close(ch)
				}
			default:
		}
	}
}