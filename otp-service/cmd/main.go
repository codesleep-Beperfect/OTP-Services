package main

import (
	"otp-service/config"
	"otp-service/internal/db"
	"otp-service/internal/router"
)

func main() {
	config.Load()
	db.Init()

	r := router.SetupRouter()
	r.Run(":" + config.App.Server.Port)
}