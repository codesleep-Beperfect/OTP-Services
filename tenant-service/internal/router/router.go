package router

import (
	"tenant-service/internal/db"
	"tenant-service/internal/handler"
	"tenant-service/internal/repository"
	"tenant-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// dependencies
	repo := repository.NewTenantRepo(db.DB)
	service := service.NewTenantService(repo)
	handler := handler.NewTenantHandler(service)

	// routes
	r.POST("/v1/tenant/register", handler.Register)
	r.GET("/v1/tenant/validate", handler.Validate)

	return r
}