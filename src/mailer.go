package mailer

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"mailer/config"
	"net/mail"
	"net/smtp"
	"strings"
	"sync"
)

type API struct {
	sync.Mutex
}

type Request struct {
	XMLName xml.Name `xml:"xml,omitempty"`
	Body    struct {
		To      mail.Address `xml:"to,omitempty"`
		From    mail.Address `xml:"from,omitempty"`
		Subject string       `xml:"subject,omitempty"`
		Msg     string       `xml:"msg,omitempty"`
	} `xml:"body,omitempty"`
}

func (api *API) Mail(request Request) bool {
	message := []string{}
	header := make(map[string]string)
	header["To"] = request.Body.To.String()
	header["From"] = request.Body.From.String()
	header["Subject"] = request.Body.Subject
	header["MIME-Version"] = `1.0`
	header["Content-Type"] = `text/html;charset="utf-8"`
	header["Content-Transfer-Encoding"] = `base64`
	for k, v := range header {
		message = append(message, fmt.Sprintf("%s: %s", k, v))
	}
	message = append(message, base64.StdEncoding.EncodeToString([]byte(request.Body.Msg)))
	auth := smtp.PlainAuth("", config.MailUser, config.MailPass, config.MailHost)
	addr := config.MailHost + ":" + config.MailPort
	err := smtp.SendMail(addr, auth, request.Body.From.Address, []string{request.Body.To.Address}, []byte(strings.Join(message, "\r\n")))
	if err != nil {
		return false
	}
	return true
}
