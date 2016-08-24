package main

import
(
	"os"
	"fmt"
	"runtime"
	"strconv"
	"code.google.com/p/log4go"	
)


type LogItem struct {
	Message string
}

var LogItems []LogItem

func saveLogs() {
 	logFile := log4go.NewFileLogWriter("stack.log", false)
    logFile.SetFormat("%d %t - %M (%S)")
    logFile.SetRotate(false)
    logFile.SetRotateSize(0)
    logFile.SetRotateLines(0)
    logFile.SetRotateDaily(true)

	logStack := make(log4go.Logger)
	logStack.AddFilter("file", log4go.DEBUG, logFile)
	for i := range LogItems {
		fmt.Println(LogItems[i].Message)
		logStack.Info(LogItems[i].Message)
	}
}


func goDetails(done chan bool) {
	i := 0
	for {
		var message string
		stackBuf := make([]byte,1024)
		stack := runtime.Stack(stackBuf, false)
		stack++	
		_, callerFile, callerLine, ok := runtime.Caller(0)
		message = "Goroutine from " + string(callerLine) + " " + string(callerFile) + " stack:" + 	string(stackBuf)		
		openGoroutines := runtime.NumGoroutine()

		if (ok == true) {
			message = message + callerFile
		}

		message = message + strconv.FormatInt(int64(openGoroutines),10) + " goroutines active"

		li := LogItem{ Message: message}

		LogItems = append(LogItems,li)
		if i == 20 {
			done <- true
			break
		}

		i++
	}
}


func main() {
	done := make(chan bool)	

	go goDetails(done)
	for i:= 0; i < 10; i++ {
		go goDetails(done)
	}

	for {
		select {
			case d := <-done:
				if d == true {
					saveLogs()
					os.Exit(1)
				}
		}
	}

}