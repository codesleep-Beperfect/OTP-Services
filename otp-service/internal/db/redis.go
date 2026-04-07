package db

import (
	"otp-service/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: config.App.Redis.Addr,
	})
}