package main

import (
	"io"

	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var errorlog *os.File

var logger *log.Logger

func init() {
	logdir := "/var/log/sumapatch"
	logfile := filepath.Join(logdir, "jobchecker.log")

	if err := os.MkdirAll(logdir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	errorlog, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger = log.New()
	logger.SetOutput(io.MultiWriter(os.Stdout, errorlog))
	formatter := &log.TextFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		FullTimestamp:   true,
	}

	logger.SetFormatter(formatter)
	logger.Infof("Logging to: %v", logfile)
	/* mw := io.MultiWriter(os.Stdout, errorlog)
	log.SetOutput(mw)
	log.Printf("Logging to: %s\n", logfile) */
}
