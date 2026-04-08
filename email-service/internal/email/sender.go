package email

import (
	"fmt"
	"net/smtp"
	"email-service/config"
)

func Send(to, otp string) error {
	auth := smtp.PlainAuth("", config.App.Email.From, config.App.Email.Password, config.App.Email.SMTPHost)

	msg := []byte(fmt.Sprintf("Subject: OTP Verification\n\nYour OTP is: %s", otp))

	addr := fmt.Sprintf("%s:%s", config.App.Email.SMTPHost, config.App.Email.SMTPPort)

	return smtp.SendMail(addr, auth, config.App.Email.From, []string{to}, msg)
}