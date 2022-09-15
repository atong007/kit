package mail

import (
	"fmt"
	"net/smtp"
)

type SmtpSender struct {
	host, port, defaultMail, passwd string
}

func NewSmtpSender(host string, port string, defaultMail string, passwd string) *SmtpSender {
	return &SmtpSender{host: host, port: port, defaultMail: defaultMail, passwd: passwd}
}

func (s *SmtpSender) SendTo(to []string, title, content string) error {
	address := s.host + ":" + s.port

	subject := "Subject: " + title + "\r\n"
	body := content + "\r\n"
	msg := []byte(subject + "\r\n" + body)

	auth := smtp.PlainAuth("", s.defaultMail, s.passwd, s.host)

	err := smtp.SendMail(address, auth, s.defaultMail, to, msg)
	if err != nil {
		return fmt.Errorf("error sending mail to %s: %v", to, err)
	}
	return nil
}
