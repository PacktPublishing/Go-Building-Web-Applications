package main

import(
"os"
)

func main() {
	
	file, _ := os.Create("/defer.txt")

	defer file.Close()
	
	for {

		// a bunch of code that extends for many lines

	}
	

}
