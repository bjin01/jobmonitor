package pkg_updates

import (
	"io"

	"os"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var errorlog *os.File

var logger *log.Logger

func Setup_Logger(log_file string) {

	errorlog, err := os.OpenFile(log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger = log.New()
	myformatter := &prefixed.TextFormatter{
		FullTimestamp:   false,
		ForceColors:     true,
		TimestampFormat: "2006/01/02 15:04:05",
	}
	logger.SetFormatter(myformatter)

	logger.SetOutput(io.MultiWriter(os.Stdout, errorlog))
	/* formatter := &log.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		PrettyPrint:     true,

		//FullTimestamp:   true,
	}

	logger.SetFormatter(formatter) */
	logger.Infof("Logging to: %v", log_file)
	/* mw := io.MultiWriter(os.Stdout, errorlog)
	log.SetOutput(mw)
	log.Printf("Logging to: %s\n", logfile) */
}

func SetLoggerLevel(level string) {
	switch level {
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
}
