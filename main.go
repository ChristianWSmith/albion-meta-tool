package main

import (
	"fmt"
	"os"
)

var config Config = defaultConfig()

func crash(message string, err error) {
	logError(message, err)
	os.Exit(1)
}

func main() {
	var err error

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

	logInfo(fmt.Sprintf("Config: %v", config), nil)
	// Your application logic here

	eventMonitor()
}
