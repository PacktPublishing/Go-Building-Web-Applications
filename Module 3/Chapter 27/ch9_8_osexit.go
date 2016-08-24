package main

import
(
	"fmt"
	"os"
)

func processNumber(un int) {

	if un < 1 || un > 4 {
		fmt.Println("Now you've done it!")
		os.Exit(1)
	}else {
		fmt.Println("Good, you can read simple instructions.")
	}
}

func main() {
	userNum := 0
	fmt.Println("Enter a number between 1 and 4.")
	_,err := fmt.Scanf("%d",&userNum)
		if err != nil {}
	
	processNumber(userNum)
}