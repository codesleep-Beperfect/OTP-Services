package model

type OTPData struct {
	Hash        string `json:"hash"`
	ResendCount int    `json:"resend_count"`
}