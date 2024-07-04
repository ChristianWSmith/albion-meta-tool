package main

import (
	"fmt"
	"os"
)

func main() {

	var config, err = getConfig()
	if err != nil {
		fmt.Println("Program failed while getting config:", err)
		os.Exit(1)
	}

	fmt.Printf("Config: %+v\n", config)
	// Your application logic here
}
