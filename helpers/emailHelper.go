package helper

import (
    "net/smtp"
    "os"
)

// sendEmail sends an email using SMTP
func SendEmail(to []string, subject string, body string) error {
    auth := smtp.PlainAuth(
        "",
        os.Getenv("FROM_EMAIL"),
        os.Getenv("FROM_EMAIL_PASSWORD"),
        os.Getenv("FROM_EMAIL_SMTP"),
    )

    message := "Subject: " + subject + "\n" + body
    return smtp.SendMail(
        os.Getenv("SMTP_ADDR"),
        auth,
        os.Getenv("FROM_EMAIL"),
        to,
        []byte(message),
    )
}
