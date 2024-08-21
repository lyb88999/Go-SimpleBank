package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddress   = "smtp.qq.com"
	smtpServerAddress = "smtp.qq.com:587"
)

type EmailSender interface {
	SendEmail(subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string) error
}

type QQEmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func (sender QQEmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	for _, file := range attachFiles {
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("failed to attach file: %s", err)
		}
	}
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}

func NewQQEmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &QQEmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}
