package lib

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	ApiURL       string `mapstructure:"api_url"`
	ApiKey       string `mapstructure:"api_key"`
	DashboardURl string `mapstructure:"dashboard_url"`
}

func CreateConfigFile() error {
	if err := os.MkdirAll(DefaultConfigFilePath(), 0777); err != nil {
		return err
	}

	if _, err := os.Create(filepath.Join(DefaultConfigFilePath(), ".config.json")); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal("unable to generate file config file: ", err)
		}
	}
	viper.AddConfigPath(DefaultConfigFilePath())
	viper.SetConfigName(".config")
	viper.SetConfigType("json")
	if environment := os.Getenv("ENVIRONMENT"); environment == "development" {
		viper.Set("api_url", "http://localhost:4000")
		viper.Set("dashboard_url", "http://localhost:3000")

	} else {
		viper.Set("api_url", DefaultBaseApiURL)

	}
	return viper.WriteConfig()

}

func LoadConfig() (config Config, err error) {

	if err = viper.ReadInConfig(); err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)

	return config, err
}
