package main

import (
	"io"

	"os"
	"path/filepath"

	"github.com/bjin01/jobmonitor/auth"
	"github.com/bjin01/jobmonitor/email"
	"github.com/bjin01/jobmonitor/pkg_updates"
	"github.com/bjin01/jobmonitor/saltapi"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = logrus.New()

func Setup_Logger(logfile string) {
	if logfile == "" {
		logdir := "/var/log/sumapatch"
		logfile = filepath.Join(logdir, "jobchecker.log")

		if err := os.MkdirAll(logdir, 0755); err != nil {
			logrus.Fatalf("Failed to create log directory: %v", err)
		}
	}

	errorlog, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	myformatter := &prefixed.TextFormatter{
		FullTimestamp:   false,
		ForceColors:     true,
		TimestampFormat: "2006/01/02 15:04:05",
	}
	logger.SetFormatter(myformatter)

	logger.SetOutput(io.MultiWriter(os.Stdout, errorlog))
	/* formatter := &log.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		//FullTimestamp:   true,
	}


	logger.SetFormatter(formatter) */
	lumberjackLogrotate := &lumberjack.Logger{
		Filename:   logfile, // File to write logs to.
		MaxSize:    1,       // Maximum file size before rotation, in megabytes.
		MaxBackups: 1,       // Maximum number of old log files to keep.
		MaxAge:     1,       // Maximum number of days to retain old log files.
		Compress:   true,    // Whether to compress the rotated log files.
	}

	logger.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogrotate))

	logger.Infof("Set logging to file: %s", errorlog.Name())

	// Pass the logger to the subpackages
	pkg_updates.SetLogger(logger)
	saltapi.SetLogger(logger)
	email.SetLogger(logger)
	auth.SetLogger(logger)

}
