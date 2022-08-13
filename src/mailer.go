package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

type Config struct {
	MailHost string
	MailPort string
	MailUser string
	MailPass string
}

type API struct {
	Buffer  *bytes.Buffer
	Config  Config
	Header  []string
	Content []string
	Body    struct {
		To      mail.Address
		From    mail.Address
		Subject interface{}
		Message interface{}
	}
}

func (api *API) SetHeaders(from, to mail.Address, subject, message interface{}) {
	api.Buffer = bytes.NewBuffer(nil)
	api.Body.From = from
	api.Body.To = to
	api.Body.Subject = subject
	api.Body.Message = message
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "To", api.Body.To.String()))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "From", api.Body.From.String()))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "Subject", api.Body.Subject))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "Mime-Version", "1.0"))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset="utf-8"`))
	api.Header = append(api.Header, fmt.Sprintf("%s", ""))
	api.Header = append(api.Header, fmt.Sprintf("%s", ""))
	api.Buffer.WriteString(strings.Join(api.Header, "\r\n"))
	api.Content = append(api.Content, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset="utf-8"`))
	api.Content = append(api.Content, fmt.Sprintf("%s", ""))
	api.Content = append(api.Content, fmt.Sprintf("%s", api.Body.Message))
	api.Buffer.WriteString(strings.Join(api.Content, "\r\n"))
	return
}

func (api *API) AddHeader(key, value string) {
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", key, value))
}

func (api *API) Send() bool {
	addr := api.Config.MailHost + ":" + api.Config.MailPort
	host, _, _ := net.SplitHostPort(addr)
	auth := smtp.PlainAuth("", api.Config.MailUser, api.Config.MailPass, host)
	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: host}
	c, err := smtp.Dial(addr)
	if err != nil {
		log.Println(err)
		return false
	}
	c.StartTLS(tlsconfig)
	if err = c.Auth(auth); err != nil {
		log.Println(err)
		return false
	}
	if err = c.Mail(api.Body.From.Address); err != nil {
		log.Println(err)
		return false
	}
	if err = c.Rcpt(api.Body.To.Address); err != nil {
		log.Println(err)
		return false
	}
	w, err := c.Data()
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = w.Write(api.Buffer.Bytes())
	if err != nil {
		log.Println(err)
		return false
	}
	err = w.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	c.Quit()
	return true
}
