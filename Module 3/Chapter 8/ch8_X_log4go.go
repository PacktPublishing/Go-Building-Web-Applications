package main

import
(
	logger "code.google.com/p/log4go"
)
func main() {
	logMech := make(logger.Logger);
	logMech.AddFilter("stdout", logger.DEBUG, logger.NewConsoleLogWriter())

	flw := logger.NewFileLogWriter("log_manager.log", false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(true)
	flw.SetRotateSize(256)
	flw.SetRotateLines(20)
	flw.SetRotateDaily(true)
	logMech.AddFilter("file", logger.FINE, flw)


	logMech.Trace("Received message: %s)", "All is well")
	logMech.Info("Message received: ", "debug!")
	logMech.Error("Oh no!","Something Broke")
}