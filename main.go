package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Database   string `toml:"database"`
}

func defaultConfig() Config {
    return Config{
        Database:   "db.sqlite",
    }
}

func loadConfig(path string) (Config, error) {
    var config Config
    if _, err := toml.DecodeFile(path, &config); err != nil {
        return Config{}, err
    }
    return config, nil
}

func saveConfig(path string, config Config) error {
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

func main() {
    var configPath string
    flag.StringVar(&configPath, "config", "amt.toml", "path to config file")
    flag.Parse()

    var config Config
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        fmt.Println("Config file not found, generating default config...")
        config = defaultConfig()
        if err := saveConfig(configPath, config); err != nil {
            fmt.Println("Error saving default config:", err)
            return
        }
    } else {
        fmt.Println("Loading config file...")
        var err error
        config, err = loadConfig(configPath)
        if err != nil {
            fmt.Println("Error loading config file:", err)
            return
        }
    }

    fmt.Printf("Config: %+v\n", config)
    // Your application logic here
}

