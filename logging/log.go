package logging

import (
	"log"
	"os"
)

const (
	GinLog = "gin.log"
	AppLog = "app.log"
)

// todo wrap logger usage with an interface

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	f, err := os.OpenFile(AppLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(f, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
