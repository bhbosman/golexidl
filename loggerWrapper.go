package main

import "log"

type loggerWrapper struct {
	applicationLogger *log.Logger
}

func (self *loggerWrapper) Printf(format string, v ...interface{}) {
	self.applicationLogger.Printf(format, v...)
}
