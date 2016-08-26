package main

import (
	"log"
	"os"
)

func main() {
	logFile, _ := os.OpenFile("C:\\wamp\\www\\test.log", os.O_RDWR, 0755)

	log.SetOutput(logFile)
	log.Println("Sending an entry to log!")

	logFile.Close()
}
