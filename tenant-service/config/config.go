package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	MySQL struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
}

var App Config

func Load() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&App); err != nil {
		log.Fatal(err)
	}
}