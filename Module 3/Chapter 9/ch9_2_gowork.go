package main

import
(
	"log"
	"os"
	"strconv"
)

const totalGoroutines = 5

type Worker struct {
	wLog *log.Logger
	Name string
	FileName string
	File *os.File
}

func main() {
	done := make(chan bool)

	for i:=0; i< totalGoroutines; i++ {

		myWorker := Worker{}
		myWorker.Name = "Goroutine " + strconv.FormatInt(int64(i),10) + " " 
		myWorker.FileName = "C:\\wamp\\www\\log_"+strconv.FormatInt(int64(i),10) + ".log" 
		tmpFile,_ :=   os.OpenFile(myWorker.FileName, os.O_CREATE, 0755)
		myWorker.File = tmpFile
		myWorker.wLog = log.New(myWorker.File, myWorker.Name, 1)
		go func(w *Worker) {

				w.wLog.Print("Hmm")

				done <- true
		}(&myWorker)
	}	

	log.Println("...")

	<- done
}