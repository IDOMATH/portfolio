package util

import (
	"errors"
	"net/smtp"
	"os"
)

func SendEmail(email, message string) error {

	from := os.Getenv("EMAIL")
	if from == "" {
		return errors.New("could not get email from .env")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		return errors.New("could not get password from .env")
	}

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
