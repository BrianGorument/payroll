package config

import (
	"log"

	"github.com/spf13/viper"
)

// LoadConfig loads configuration from .env file using Viper
func LoadConfig() error {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return err
	}

	return nil
}
