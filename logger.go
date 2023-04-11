package main

import (
	"io"
	"log"
	"os"
)

var errorlog *os.File

//var logger *log.Logger

func init() {
	logfile := "/var/log/sumapatch/jobchecker.log"
	errorlog, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	mw := io.MultiWriter(os.Stdout, errorlog)
	log.SetOutput(mw)
	log.Printf("Logging to: %s\n", logfile)
}
