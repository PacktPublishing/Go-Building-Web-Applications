package main

import
(
"fmt"
"net"
"os"
"bufio"
"strings"

)

type Message struct {
	message string
	user string
}

var recvBuffer [140]byte

func listen(conn net.Conn) {
	for {

	    messBuff := make([]byte,1024)
			n, err := conn.Read(messBuff)
			if err != nil {

			}
			message := string(messBuff[:n])
			message = message[0:]

			fmt.Println(strings.TrimSpace(message))
			fmt.Print("> ")
	}

}

func talk(conn net.Conn, mS chan Message) {

      for {
			command := bufio.NewReader(os.Stdin)
				fmt.Print("> ")      	
                line, err := command.ReadString('\n')
                
                line = strings.TrimRight(line, " \t\r\n")
				_, err = conn.Write([]byte(line))                       
                if err != nil {
                        conn.Close()
                        break

                }
			doNothing(command)	
        }	

}

func doNothing(bf *bufio.Reader) {


}
func main() {

	messageServer := make(chan Message)

	userName := os.Args[1]

	fmt.Println("Connecting to host as",userName)

	clientClosed := make(chan bool)

	conn, err := net.Dial("tcp","127.0.0.1:9000")
	if err != nil {
		fmt.Println("Could not connect to server!")
	}
	conn.Write([]byte(userName)) 
	introBuff := make([]byte,1024)	  
	n, err := conn.Read(introBuff)
	if err != nil {

	}
	message := string(introBuff[:n])	
	fmt.Println(message)

	go talk(conn,messageServer)
	go listen(conn)

	<- clientClosed
}
