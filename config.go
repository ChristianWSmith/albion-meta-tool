package main

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Database             string        `toml:"database"`
	KillEventUrl         string        `toml:"albion_event_url"`
	PriceUrl             string        `toml:"albion_online_data_url"`
	ItemNamesUrl         string        `toml:"item_names_url"`
	PriceLocations       []string      `toml:"price_locations"`
	PriceStaleThreshold  time.Duration `toml:"price_stale_threshold"`
	EventStaleThreshold  time.Duration `toml:"event_stale_threshold"`
	EventCleanupInterval time.Duration `toml:"event_cleanup_interval"`
	LogFile              string        `toml:"log_file"`
	LogLevel             logrus.Level  `toml:"log_level"`
	PollEvents           bool          `toml:"poll_events"`
	Port                 int           `toml:"port"`
}

func defaultConfig() Config {
	return Config{
		Database:             "amt.sqlite",
		KillEventUrl:         "https://gameinfo.albiononline.com/api/gameinfo/events",
		PriceUrl:             "https://old.west.albion-online-data.com/api/v2/stats/History",
		ItemNamesUrl:         "https://raw.githubusercontent.com/ao-data/ao-bin-dumps/master/formatted/items.txt",
		PriceLocations:       []string{"Lymhurst", "Thetford", "FortSterling", "Martlock", "Bridgewatch"},
		PriceStaleThreshold:  time.Duration(7*24) * time.Hour,
		EventStaleThreshold:  time.Duration(7*24) * time.Hour,
		EventCleanupInterval: time.Duration(24) * time.Hour,
		LogFile:              "amt.log",
		LogLevel:             logrus.PanicLevel,
		PollEvents:           true,
		Port:                 8080,
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
		log.Info("Config file not found, generating default config")
		config = defaultConfig()
		if err := saveConfigFile(configPath, config); err != nil {
			log.Error("Failed saving default config", err)
			return config, err
		}
	} else {
		log.Info("Loading config file")
		var err error
		config, err = loadConfigFile(configPath)
		if err != nil {
			log.Error("Failed loading config file", err)
			return config, err
		}
	}
	return config, nil
}
