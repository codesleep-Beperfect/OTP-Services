package main

import (
	"email-service/config"
	"email-service/internal/kafka"
)

func main() {
	config.Load()
	kafka.StartConsumer()
}