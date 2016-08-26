package main

import
(
	"fmt"
	"net"
	"io"
	"os"
	"strconv"
	"encoding/json"
)

var backupFolder = "C:\\wamp\\www\\backup\\"

type Message struct {
	Hash string "json:hash"
	Action string "json:action"
	Location string "json:location"
	Name string "json:name"	
	Version int "json:version"
}

func backup (location string, name string, version int) {

	newFileName := backupFolder + name + "." + strconv.FormatInt(int64(version),10)
	fmt.Println(newFileName)
	org,_ := os.Open(location)
	defer org.Close()
	cpy,_ := os.Create(newFileName)
	defer cpy.Close()
	io.Copy(cpy,org)
}

func listen(conn net.Conn) {
	for {

	    messBuff := make([]byte,1024)
		n, err := conn.Read(messBuff)
		if err != nil {

		}
		message := string(messBuff[:n])
		message = message[0:]

		resultMessage := Message{}
		json.Unmarshal(messBuff[:n],&resultMessage)
		
		if resultMessage.Action == "MODIFY" {
			fmt.Println("Back up file",resultMessage.Location)
			newVersion := resultMessage.Version + 1
			backup(resultMessage.Location,resultMessage.Name,newVersion)
		}
		
	}

}

func main() {
	endBackup := make(chan bool)
	conn, err := net.Dial("tcp","127.0.0.1:9000")
	if err != nil {
		fmt.Println("Could not connect to File Listener!")
	}
	go listen(conn)

	<- endBackup
}