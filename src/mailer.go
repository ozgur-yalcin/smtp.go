package mailer

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"mailer/config"
	"net/mail"
	"net/smtp"
	"strings"
)

type MailData struct {
	XMLName xml.Name `xml:"xml,omitempty"`
	Body    struct {
		To      mail.Address `xml:"to,omitempty"`
		From    mail.Address `xml:"from,omitempty"`
		Subject string       `xml:"subject,omitempty"`
		Msg     string       `xml:"msg,omitempty"`
	} `xml:"body,omitempty"`
}

func Mail(xmlrequest MailData) bool {
	message := []string{}
	header := make(map[string]string)
	header["To"] = xmlrequest.Body.To.String()
	header["From"] = xmlrequest.Body.From.String()
	header["Subject"] = xmlrequest.Body.Subject
	header["MIME-Version"] = `1.0`
	header["Content-Type"] = `text/html;charset="utf-8"`
	header["Content-Transfer-Encoding"] = `base64`
	for k, v := range header {
		message = append(message, fmt.Sprintf("%s: %s", k, v))
	}
	message = append(message, base64.StdEncoding.EncodeToString([]byte(xmlrequest.Body.Msg)))
	auth := smtp.PlainAuth("", config.MailUser, config.MailPass, config.MailHost)
	addr := config.MailHost + ":" + config.MailPort
	err := smtp.SendMail(addr, auth, xmlrequest.Body.From.Address, []string{xmlrequest.Body.To.Address}, []byte(strings.Join(message, "\r\n")))
	if err != nil {
		return false
	}
	return true
}
