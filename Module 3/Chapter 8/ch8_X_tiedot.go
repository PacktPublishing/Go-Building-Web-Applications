package main

import
(
	"fmt"
)

func main() {

	Collection := Collection{
		Name: ''
	}
	
	data, err := http.Get("http://localhost:8080/all")
	if (err != nil) {
		fmt.Println("Error accessing tiedot")
	}
	collections,_ = json.Unmarshal(data,&Collection)
}