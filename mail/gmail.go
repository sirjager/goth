package mail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"

	"github.com/sirjager/goth/vo"
)

// _GmailSender is a struct that represents a Gmail email sender.
type _GmailSender struct {
	address     string
	plainAuth   smtp.Auth
	senderName  string
	senderEmail string
	// log email on console, useful when testing, unsafe in production
	logMail bool
}

const (
	// _GmailSMTPHost is the SMTP host for Gmail.
	_GmailSMTPHost = "smtp.gmail.com"
	// _GmailSMTPPort is the SMTP port for Gmail.
	_GmailSMTPPort = "587"
)

// NewGmailSender creates a new _GmailSender instance with the provided GmailSMTP configuration.
// If [logMails] is provided, then emails will be logged in console and will not be sent to email.
// [logMails] option is only supposed to be used in local and testing environemt, but never in production.
func NewGmailSender(config Config, logMails ...bool) (Sender, error) {
	// Validate the SMTP user email.
	email, err := vo.NewEmail(config.SMTPUser)
	if err != nil {
		return nil, err
	}
	return &_GmailSender{
		senderEmail: email.Value(),
		senderName:  config.SMTPSender,
		address:     _GmailSMTPHost + ":" + _GmailSMTPPort,
		plainAuth:   smtp.PlainAuth("", email.Value(), config.SMTPPass, _GmailSMTPHost),
		logMail:     len(logMails) == 1 && logMails[0],
	}, nil
}

// SendMail sends an email using the _GmailSender.
func (gmail *_GmailSender) SendMail(mail Mail) error {
	// Validate the recipient email addresses.
	for _, r := range mail.To {
		_, err := vo.NewEmail(r)
		if err != nil {
			return err
		}
	}

	// Validate the Bcc email addresses.
	for _, r := range mail.Bcc {
		_, err := vo.NewEmail(r)
		if err != nil {
			return err
		}
	}

	// Validate the Cc email addresses.
	for _, r := range mail.Cc {
		_, err := vo.NewEmail(r)
		if err != nil {
			return err
		}
	}

	// Create a new email instance.
	email := email.NewEmail()
	email.From = fmt.Sprintf("%s <%s>", gmail.senderName, gmail.senderEmail)
	email.To = mail.To
	email.Cc = mail.Cc
	email.Bcc = mail.Bcc
	email.Subject = mail.Subject
	email.HTML = []byte(mail.Body)

	// Attach files to the email.
	for _, f := range mail.Files {
		if _, err := email.AttachFile(f); err != nil {
			return fmt.Errorf("failed to attach file: %s :%w", f, err)
		}
	}

	if gmail.logMail {
		fmt.Println("--------------------------------------------------------------------")
		fmt.Println("to: ", mail.To)
		fmt.Println("cc: ", strings.Join(mail.Cc, ","))
		fmt.Println("bcc: ", strings.Join(mail.Bcc, ","))
		fmt.Println("subject: ", mail.Subject)
		fmt.Println("body: ", mail.Body)
		fmt.Println("--------------------------------------------------------------------")
		return nil
	}

	// Send the email using the Gmail SMTP server and authentication.
	return email.Send(gmail.address, gmail.plainAuth)
}
