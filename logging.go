package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func initLogging() error {
	log.SetFormatter(&logrus.JSONFormatter{})
	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error("Failed to open log file: ", config.LogFile)
		return err
	}
	log.SetOutput(logFile)
	log.SetLevel(config.LogLevel)
	return nil
}
