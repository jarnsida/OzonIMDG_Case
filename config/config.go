package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	config Config
	once   sync.Once
)

type Config struct {
	LogLevel         string  `envconfig:"LOG_LEVEL"`
	MaxMemory        string  `envconfig:"MAX_MEMORY"`
	FilePath         string  `envconfig:"FILE_PATH"`
	IMDBPort         string  `envconfig:"PORT"`
	ConnCloseTimeout float64 `envconfig:"connCloseTO"`
}

// Чтение конфигурации из среды. Once.
func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Чтение конфигурации из среды:", string(configBytes))
	})
	return &config
}
