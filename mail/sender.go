package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

const (
	GmailHost        = "smtp.gmail.com"
	GmailHostAddress = "smtp.gmail.com:587"
)

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.Cc = cc
	e.To = to
	e.Bcc = bcc
	e.HTML = []byte(content)
	for _, attachFile := range attachFiles {
		_, err := e.AttachFile(attachFile)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", attachFile, err)
		}
	}
	fmt.Print(2345678)
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, GmailHost)
	return e.Send(GmailHostAddress, smtpAuth)
}
