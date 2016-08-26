package main

import (
	"runtime"
	"fmt"
)

func listThreads()(int) {
	threads := runtime.GOMAXPROCS(0)
	return threads
}

func main() {

	fmt.Printf("%d thread(s) available to Go.",listThreads());

}