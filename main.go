package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var config Config = defaultConfig()
var log = logrus.New()

func crash(message string, err error) {
	log.Error(message, err)
	os.Exit(1)
}

func main() {
	var err error
	log.SetLevel(config.LogLevel)

	config, err = getConfig()
	if err != nil {
		crash("Failed while getting config", err)
	}

	err = initLogging()
	if err != nil {
		crash("Failed to initialize logging", err)
	}

	err = initDatabase()
	if err != nil {
		crash("Failed to initialize database", err)
	}

	log.Info("Config: ", config)
	// Your application logic here

	eventMonitor()
}
