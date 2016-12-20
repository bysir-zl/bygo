package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type Mail struct {
	from, pass, smtpHost, smtpPort string
}

func NewMail(from, pass, smtpHostP string) *Mail {
	smtpHost := smtpHostP
	port := "465"
	if strings.Contains(smtpHostP, ":") {
		hAp := strings.Split(smtpHostP, ":")
		smtpHost = hAp[0]
		port = hAp[1]
	}
	return &Mail{
		from:    from,
		pass:    pass,
		smtpHost:smtpHost,
		smtpPort:port,
	}
}

func (p *Mail) Send(subject, body string, to []string) error {
	auth := smtp.PlainAuth("", p.from, p.pass, p.smtpHost)

	var conn net.Conn
	var err error

	// 判断是否用了ssl
	if p.smtpPort == "465" {
		conn, err = tls.Dial("tcp", p.smtpHost + ":" + p.smtpPort, nil)
	} else {
		conn, err = net.DialTimeout("tcp", p.smtpHost + ":" + p.smtpPort, 10 * time.Second)
	}

	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, p.smtpHost)
	if err != nil {
		return err
	}
	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(p.from); err != nil {
		return err
	}

	contentType := "Content-Type: text/plain;charset=UTF-8"
	msgStr := fmt.Sprint(
		"To:", strings.Join(to, ";"),
		"\r\nFrom:", p.from,
		"\r\nSubject:", subject,
		"\r\n", contentType,
		"\r\n\r\n", body,
	)
	msg := []byte(msgStr)
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	if err != nil {
		return err
	}
	writer.Close()
	return nil
}
