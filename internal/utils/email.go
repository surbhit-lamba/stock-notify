package utils

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"os"
	"stock-notify/pkg/log"
	"text/template"
)

func SendEmailWithHTMLTemplate(ctx context.Context, from string, to []string, templatePathFromRoot string, subject string) {
	t, err := template.ParseFiles(templatePathFromRoot)
	if err != nil {
		log.ErrorfWithContext(ctx, "[SendEmailWithHTMLTemplate] Parsing Error - ", err.Error())
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: "+subject+"\n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Yolo Bhai",
		Message: "This is a test stock message in a HTML template",
	})
	if err != nil {
		log.ErrorfWithContext(ctx, "[SendEmailWithHTMLTemplate] Error - ", err.Error())
	}

	SendEmail(ctx, from, to, body)
}

func SendEmail(ctx context.Context, from string, to []string, body bytes.Buffer) {
	// smtp server configuration.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Authentication.
	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.ErrorfWithContext(ctx, "[SendEmail] could not send email - ", err.Error())
		return
	}
	log.InfofWithContext(ctx, "[SendEmail] sent mail - ", from, to, string(body.Bytes()))
}
