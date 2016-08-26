package main

import(
	_"runtime"
	"fmt"

)

func showNumber(num int) {
	fmt.Println(num)
}

func main() {
	iterations := 10
	
	for i := 0; i<=iterations; i++ {

		go showNumber(i)

	}
	//runtime.Gosched()
	fmt.Println("Goodbye!")

}