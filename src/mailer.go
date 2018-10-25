package mailer

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"path"
	"strings"
	"sync"

	"github.com/OzqurYalcin/mailer/config"
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
		Files   []string     `xml:"files,omitempty"`
	} `xml:"body,omitempty"`
}

func (api *API) Boundary() string {
	var buf [30]byte
	io.ReadFull(rand.Reader, buf[:])
	return string(buf[:])
}

func (api *API) Mail(request *Request) bool {
	message := []string{}
	message = append(message, fmt.Sprintf("%s: %s", "To", request.Body.To.String()))
	message = append(message, fmt.Sprintf("%s: %s", "From", request.Body.From.String()))
	message = append(message, fmt.Sprintf("%s: %s", "Subject", request.Body.Subject))
	message = append(message, fmt.Sprintf("%s: %s", "MIME-Version", `1.0`))
	if len(request.Body.Files) > 0 {
		boundary := api.Boundary()
		message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `multipart/mixed;boundary=`+boundary))
		message = append(message, fmt.Sprintf("--%s", boundary))
		message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset="utf-8"`))
		message = append(message, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
		message = append(message, base64.StdEncoding.EncodeToString([]byte(request.Body.Msg)))
		message = append(message, fmt.Sprintf("--%s", boundary))
		for _, file := range request.Body.Files {
			content, err := ioutil.ReadFile(file)
			if err == nil {
				message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `application/octet-stream`))
				message = append(message, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
				message = append(message, fmt.Sprintf("%s: %s", "Content-Disposition", `attachment; filename=`+path.Base(file)))
				message = append(message, base64.StdEncoding.EncodeToString(content))
				message = append(message, fmt.Sprintf("--%s", boundary))
			} else {
				return false
			}
		}
	} else {
		message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset="utf-8"`))
		message = append(message, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
		message = append(message, base64.StdEncoding.EncodeToString([]byte(request.Body.Msg)))
	}
	auth := smtp.PlainAuth("", config.MailUser, config.MailPass, config.MailHost)
	addr := config.MailHost + ":" + config.MailPort
	err := smtp.SendMail(addr, auth, request.Body.From.Address, []string{request.Body.To.Address}, []byte(strings.Join(message, "\r\n")))
	if err != nil {
		return false
	}
	return true
}
