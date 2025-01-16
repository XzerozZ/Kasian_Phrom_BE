package utils

import (
	"bytes"
	"strconv"
	"text/template"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gopkg.in/gomail.v2"
)

func SendMail(templatePath string, user entities.User, otp string, config configs.Mail) error {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	err = t.Execute(&body, struct { 
		Username 	string 
		OTP   		string 
	}{ 
		Username: 	user.Username, 
		OTP:   		otp,
	})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Recovery Your Password")
	m.SetBody("text/html", body.String())
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(config.Host, port, config.Sender, config.Key)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}