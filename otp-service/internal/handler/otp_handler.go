package handler

import (
	
	"otp-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OTPHandler struct {
	svc *service.OTPService
}

func NewOTPHandler(s *service.OTPService) *OTPHandler {
	return &OTPHandler{svc: s}
}

type otpRequest struct {
	Identifier string `json:"identifier"`
	OTP        string `json:"otp,omitempty"`
}

func (h *OTPHandler) Send(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")
	var req otpRequest
	if err := c.BindJSON(&req); err != nil {
	c.JSON(400, gin.H{"error": "invalid request"})
	return
	}
	
	otp, err := h.svc.Send(apiKey, req.Identifier)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"otp": otp})
}

func (h *OTPHandler) Resend(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")

	var req otpRequest
	
	if err := c.BindJSON(&req); err != nil {
	c.JSON(400, gin.H{"error": "invalid request"})
	return
	}

	otp, err := h.svc.Resend(apiKey, req.Identifier)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"otp": otp})
}

func (h *OTPHandler) Verify(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")

	var req otpRequest

	if err := c.BindJSON(&req); err != nil {
	c.JSON(400, gin.H{"error": "invalid request"})
	return
	}

	ok, err := h.svc.Verify(apiKey, req.Identifier, req.OTP)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"valid": ok})
}