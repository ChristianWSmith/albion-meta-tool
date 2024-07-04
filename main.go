package main

import (
	"fmt"
	"log"
	"os"
)

func crash(message string, err error) {
	log.Fatal(message+":", err)
	os.Exit(1)
}

func main() {

	var config, err = getConfig()
	if err != nil {
		crash("Failed while getting config", err)
	}

	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		crash("Failed to open log file", err)
	}

	log.SetOutput(logFile)

	err = initDatabase(config)
	if err != nil {
		crash("Failed to initialize database", err)
	}

	fmt.Printf("Config: %+v\n", config)
	// Your application logic here
}
