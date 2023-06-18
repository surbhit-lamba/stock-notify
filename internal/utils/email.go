package utils

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"stock-notify/pkg/log"
	"text/template"
)

func SendEmailWithHTMLTemplate(ctx context.Context) {
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

	t, err := template.ParseFiles("emailtemplates/template.html")
	if err != nil {
		log.ErrorfWithContext(ctx, "could not parse email")
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
}

func SendEmail(ctx context.Context, from string, to []string, body []byte) {
	// smtp server configuration.
	smtpHost := "mail.smtp2go.com"
	smtpPort := "2525"
	username := "stock-notifiers"
	password := "hP3b8iDkmGFl39IX"

	// Authentication.
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body)
	if err != nil {
		log.ErrorfWithContext(ctx, "[SendEmail] could not send email - ", err.Error())
		return
	}
	log.InfofWithContext(ctx, "[SendEmail] sent mail - ", from, to, string(body))
}
