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
		crash("Failed while getting config: ", err)
	}

	err = initLogging()
	if err != nil {
		crash("Failed to initialize logging: ", err)
	}

	err = initDatabase()
	if err != nil {
		crash("Failed to initialize database: ", err)
	}

	log.Info("Config: ", config)
	// Your application logic here

	if config.PollEvents {
		go eventMonitor()
	}

	prices, _ := getItemPrices([]Item{
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        4,
			Enchantment: 0,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        4,
			Enchantment: 1,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        4,
			Enchantment: 2,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        4,
			Enchantment: 3,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        4,
			Enchantment: 4,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        5,
			Enchantment: 0,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        5,
			Enchantment: 1,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        5,
			Enchantment: 2,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        5,
			Enchantment: 3,
			Quality:     0,
		},
		{
			Name:        "MAIN_CURSEDSTAFF_UNDEAD",
			Tier:        5,
			Enchantment: 4,
			Quality:     0,
		},
	})
	log.Info("Prices: ", prices)

	select {}
}
