package config

import (
	"log"
	"github.com/spf13/viper"
)

type Config struct {
	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		Topic   string	`mapstructure:"topic"`
	} `mapstructure:"kafka"`
	Email struct {
		From     string	`mapstructure:"from"`
		Password string	`mapstructure:"password"`
		SMTPHost string `mapstructure:"smtp_host"`
		SMTPPort string `mapstructure:"smtp_port"`
	} `mapstructure:"email"`
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