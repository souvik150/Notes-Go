package utils

import (
	"fmt"
	initializers "github.com/souvik150/golang-fiber/config"
	"net/smtp"
)

func SendEmail(recipient string, msg string) (string, error) {
	config, _ := initializers.LoadConfig(".")
	auth := smtp.PlainAuth(
		"",
		config.Email,
		config.EmailPassword,
		"smtp.gmail.com",
	)

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		config.Email,
		[]string{recipient},
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}

	return "", nil
}
