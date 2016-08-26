package main

import
(
	"net"
	"fmt"
)

type Connection struct {

}

func (c Connection) Listen(l net.Listener) {
	for {
		conn,_ := l.Accept()
		go c.logListen(conn)
	}
}

func (c *Connection) logListen(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		fmt.Println("Log Message",string(n))
	}
}

func main() {
	serverClosed := make(chan bool)

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		fmt.Println ("Could not start server!",err)
	}

	Conn := Connection{}

	go Conn.Listen(listener)	

	<-serverClosed
}