package handler

import (
	"net/mail"

	"github.com/gin-gonic/gin"
	"tenant-service/internal/service"
)

type TenantHandler struct {
	svc *service.TenantService
}

func NewTenantHandler(s *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: s}
}

type registerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *TenantHandler) Register(c *gin.Context) {
	var req registerRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(400, gin.H{"error": "invalid email"})
		return
	}

	t, err := h.svc.Register(req.Name, req.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, t)
}

func (h *TenantHandler) Validate(c *gin.Context) {
	apiKey := c.Query("api_key")

	t, err := h.svc.Validate(apiKey)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid api key"})
		return
	}

	c.JSON(200, gin.H{
		"id": t.ID,
	})
}