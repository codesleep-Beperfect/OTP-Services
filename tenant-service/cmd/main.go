package main

import (
	"tenant-service/config"
	"tenant-service/internal/db"
	"tenant-service/internal/router"
)

func main() {
	config.Load()
	db.Init()

	r := router.SetupRouter()
	r.Run(":" + config.App.Server.Port)
}