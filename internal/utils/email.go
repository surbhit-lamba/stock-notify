package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

func SendEmail() {

	// Sender data.
	from := "surbhitla@gmail.com"
	password := "hP3b8iDkmGFl39IX"

	// Receiver email address.
	to := []string{
		"surbhitla@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "mail.smtp2go.com"
	smtpPort := "2525"

	// Authentication.
	auth := smtp.PlainAuth("stock-notifiers", "stock-notifiers", password, smtpHost)

	t, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("lag gye na ", err.Error())
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test stock notify \n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Yolo Bhai",
		Message: "This is a test stock message in a HTML template",
	})
	if err != nil {
		return
	}

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!!")
}
