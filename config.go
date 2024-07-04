package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database            string        `toml:"database"`
	KillEventUrl        string        `toml:"albion_event_url"`
	PriceUrl            string        `toml:"albion_online_data_url"`
	PriceLocations      []string      `toml:"price_locations"`
	PriceStaleThreshold time.Duration `toml:"price_stale_threshold"`
	LogFile             string        `toml:"log_file"`
	LogLevel            string        `toml:"log_level"`
	Port                int           `toml:"port"`
}

func defaultConfig() Config {
	return Config{
		Database:            "amt.sqlite",
		KillEventUrl:        "https://gameinfo.albiononline.com/api/gameinfo/events",       // ?limit={limit}&offset={offset}
		PriceUrl:            "https://old.west.albion-online-data.com/api/v2/stats/Prices", // {itemList}.json
		PriceLocations:      []string{"Lymhurst", "Thetford", "FortSterling", "Martlock", "Bridgewatch"},
		PriceStaleThreshold: time.Duration(7*24) * time.Hour,
		LogFile:             "amt.log",
		LogLevel:            "info",
		Port:                8080,
	}
}

func loadConfigFile(path string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func saveConfigFile(path string, config Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		return err
	}
	return nil
}

func getConfig() (Config, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "amt.toml", "path to config file")
	flag.Parse()

	var config Config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file not found, generating default config...")
		config = defaultConfig()
		if err := saveConfigFile(configPath, config); err != nil {
			fmt.Println("Error saving default config:", err)
			return config, err
		}
	} else {
		fmt.Println("Loading config file...")
		var err error
		config, err = loadConfigFile(configPath)
		if err != nil {
			fmt.Println("Error loading config file:", err)
			return config, err
		}
	}
	return config, nil
}
