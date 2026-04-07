package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashOTP hashes OTP with tenantID and identifier (salted hashing)
func HashOTP(tenantID, identifier, otp string) string {
	data := fmt.Sprintf("%s:%s:%s", tenantID, identifier, otp)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:])
}