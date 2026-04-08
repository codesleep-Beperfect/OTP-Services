package model

type OTPEvent struct {
	TenantID  string `json:"tenant_id"`
	Identifier string `json:"identifier"`
	OTP       string `json:"otp"`
	ExpiresAt int64  `json:"expires_at"`
}