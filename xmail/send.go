package xmail

import (
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

// Provider can be null,
// if null, will try to parse built-in provider by From email address.
func Send(e Envelope, c SendContent, password string, p *Provider) error {
	// Validate
	if p == nil {
		err := error(nil)
		p, err = TryParseProvider(e.From.Email)
		if err != nil {
			return err
		}
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if err := e.Validate(); err != nil {
		return err
	}
	if len(password) == 0 {
		return errors.Errorf("Send(): Empty password")
	}

	// Convert Envelope to gomail Message
	msg := gomail.NewMessage()
	msg.SetHeader("Subject", e.Subject)
	msg.SetHeader("From", e.From.Email)
	to := []string{}
	for _, v := range e.To {
		to = append(to, msg.FormatAddress(v.Email, v.Showname))
	}
	msg.SetHeader("To", to...)
	cc := []string{}
	for _, v := range e.Cc {
		cc = append(cc, msg.FormatAddress(v.Email, v.Showname))
	}
	msg.SetHeader("Cc", cc...)
	if c.BodyType == BodyTypeHTML {
		msg.SetBody("text/html", c.BodyString)
	} else {
		msg.SetBody("text/plain", c.BodyString)
	}
	for _, v := range c.AttachmentsPath {
		msg.Attach(v)
	}

	// Send the email
	sndAddr, sndPort, ssl, err := p.GetSendServer()
	if err != nil {
		return err
	}
	loginname, err := e.From.LoginName()
	if err != nil {
		return err
	}
	d := gomail.NewDialer(sndAddr, sndPort, loginname, password)
	if ssl {
		d.SSL = true
	}
	if err := d.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
