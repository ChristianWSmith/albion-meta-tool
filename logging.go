package main

import (
	"log"
	"os"
	"slices"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

var LOG_LEVELS = []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

var LOG_LEVELS_DEBUG = []LogLevel{LogLevelDebug}
var LOG_LEVELS_INFO = []LogLevel{LogLevelDebug, LogLevelInfo}
var LOG_LEVELS_WARN = []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn}
var LOG_LEVELS_ERROR = []LogLevel{LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError}

func initLogging() error {
	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)
	return nil
}

func logDebug(message string, err error) {
	if slices.Contains(LOG_LEVELS_DEBUG, config.LogLevel) {
		_log("DEBUG:", message, err)
	}
}

func logInfo(message string, err error) {
	if slices.Contains(LOG_LEVELS_INFO, config.LogLevel) {
		_log("INFO:", message, err)
	}
}

func logWarn(message string, err error) {
	if slices.Contains(LOG_LEVELS_WARN, config.LogLevel) {
		_log("WARNING:", message, err)
	}
}

func logError(message string, err error) {
	if slices.Contains(LOG_LEVELS_ERROR, config.LogLevel) {
		_log("ERROR:", message, err)
	}
}

func _log(prefix string, message string, err error) {
	if err != nil {
		log.Println(prefix, message+":", err)
	} else {
		log.Println(prefix, message)
	}
}
