package router

import (
	"otp-service/config"
	"otp-service/internal/client"
	"otp-service/internal/db"
	"otp-service/internal/handler"
	"otp-service/internal/repository"
	"otp-service/internal/service"
	"otp-service/internal/kafka"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// dependencies
	redisRepo := repository.NewRedisRepo(db.Client)
	tenantClient := client.NewTenantClient(config.App.TenantService.BaseURL)
	producer := kafka.NewProducer(config.App.Kafka.Brokers, config.App.Kafka.Topic)
	otpService := service.NewOTPService(redisRepo, tenantClient , producer)
	otpHandler := handler.NewOTPHandler(otpService)

	// routes
	r.POST("/v1/otp/send", otpHandler.Send)
	r.POST("/v1/otp/resend", otpHandler.Resend)
	r.POST("/v1/otp/verify", otpHandler.Verify)

	return r
}