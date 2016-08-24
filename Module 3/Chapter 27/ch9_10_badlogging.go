
package main

import
(
	"code.google.com/p/log4go"
)

type LogItem struct {
	Message string
	Function string
}

var Logs []LogItem

func SaveLogs() {
	logFile := log4go.NewFileLogWriter("errors.log",false)
    logFile.SetFormat("%d %t - %M (%S)")
    logFile.SetRotate(true)
    logFile.SetRotateSize(0)
    logFile.SetRotateLines(500)
    logFile.SetRotateDaily(false)

   	errorLog := make(log4go.Logger)
   	errorLog.AddFilter("file",log4go.DEBUG,logFile)    

   	for i:= range Logs {
   		errorLog.Info(Logs[i].Message + " in " + Logs[i].Function)
   	}

}

func registerError(block chan bool) {

	Log := LogItem{ Message:"An Error Has Occurred!", Function: "registerError()"}
	Logs = append(Logs,Log)

	block <- true
}

func separateFunction() {
	panic("Application quitting!")
}


func main() {
	block := make(chan bool)
	defer SaveLogs()
	go func(block chan bool) {

		registerError(block)

	}(block)

	separateFunction()


}