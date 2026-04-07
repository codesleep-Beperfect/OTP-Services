package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateAPIKey() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return "sk_" + hex.EncodeToString(bytes)
}