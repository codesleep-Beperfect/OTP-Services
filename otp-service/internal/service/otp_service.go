package service

import (
	"encoding/json"
	"fmt"
	"time"

	"otp-service/internal/client"
	"otp-service/internal/model"
	"otp-service/internal/repository"
	"otp-service/internal/utils"
)

type OTPService struct {
	repo   *repository.RedisRepo
	client *client.TenantClient
}

func NewOTPService(r *repository.RedisRepo, c *client.TenantClient) *OTPService {
	return &OTPService{repo: r, client: c}
}

func (s *OTPService) Send(apiKey, identifier string) (string, error) {
	tenantID, err := s.client.Validate(apiKey)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("otp:%s:%s", tenantID, identifier)
	// check if otp already existed
	exists, _ := s.repo.Exists(key)
	if exists {
		return "", fmt.Errorf("otp already generated, please wait")
	}
	// Generated OTP 
	otp := utils.GenerateOTP()
	data := model.OTPData{
		Hash:        utils.HashOTP(tenantID , identifier , otp),
		ResendCount: 0,
	}

	

	bytes, _ := json.Marshal(data)
	s.repo.Set(key, string(bytes), 5*time.Minute)

	return otp, nil
}

func (s *OTPService) Resend(apiKey, identifier string) (string, error) {
	tenantID, err := s.client.Validate(apiKey)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("otp:%s:%s", tenantID, identifier)
	// Getting OTP from Redis (if existed , otherwise error)
	val, err := s.repo.Get(key)
	if err != nil {
		return "", fmt.Errorf("otp not found")
	}

	var data model.OTPData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
	return "", err
	}

	if data.ResendCount >= 3 {
		return "", fmt.Errorf("resend limit exceeded")
	}

	newOTP := utils.GenerateOTP()
	// Hashing OTP
	data.Hash = utils.HashOTP(tenantID , identifier , newOTP)
	data.ResendCount++

	bytes, _ := json.Marshal(data)

	// overwrite existing key (best practice)
	s.repo.Set(key, string(bytes), 5*time.Minute)

	return newOTP, nil
}

func (s *OTPService) Verify(apiKey, identifier, otp string) (bool, error) {
	tenantID, err := s.client.Validate(apiKey)
	if err != nil {
		return false, err
	}

	key := fmt.Sprintf("otp:%s:%s", tenantID, identifier)

	val, err := s.repo.Get(key)
	fmt.Println(val , err)
	if err != nil {
		return false, fmt.Errorf("otp expired")
	}

	var data model.OTPData
	json.Unmarshal([]byte(val), &data)
	
	// Checking hashed otp and existing hashed otp would same or not
	if data.Hash == utils.HashOTP(tenantID , identifier , otp) {
		s.repo.Delete(key)
		return true, nil
	}

	return false, nil
}