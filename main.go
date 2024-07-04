package main

import (
	"fmt"
	"os"
)

func main() {

	var config, err = getConfig()
	if err != nil {
		fmt.Println("Failed while getting config:", err)
		os.Exit(1)
	}

	err = init_database(config)
	if err != nil {
		fmt.Println("Failed while initializing database:", err)
		os.Exit(1)
	}

	fmt.Printf("Config: %+v\n", config)
	// Your application logic here
}
