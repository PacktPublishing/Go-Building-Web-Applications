package main

import
(
	"fmt"
"github.com/couchbaselabs/go-couchbase"
)

func main() {
	
		conn, err := couchbase.Connect("http://localhost:8091")
		if err != nil {
			fmt.Println("Error:",err)
		}

		for _, pn := range conn.Info.Pools {
				fmt.Printf("Found pool:  %s -> %s\n", pn.Name, pn.URI)
		}
}