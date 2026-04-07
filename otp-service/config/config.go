package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Redis struct {
		Addr string `mapstructure:"addr"`
	} `mapstructure:"redis"`

	TenantService struct {
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"tenant_service"`
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