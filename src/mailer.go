package mailer

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		To      mail.Address  `xml:"to,omitempty"`
		From    mail.Address  `xml:"from,omitempty"`
		Subject interface{}   `xml:"subject,omitempty"`
		Msg     interface{}   `xml:"msg,omitempty"`
		Files   []interface{} `xml:"files,omitempty"`
	} `xml:"body,omitempty"`
}

func (api *API) Boundary() string {
	var buf [30]byte
	io.ReadFull(rand.Reader, buf[:])
	return fmt.Sprintf("%x", buf[:])
}

func (api *API) Mail(request *Request) bool {
	message := []string{}
	message = append(message, fmt.Sprintf("%s: %s", "To", request.Body.To.String()))
	message = append(message, fmt.Sprintf("%s: %s", "From", request.Body.From.String()))
	message = append(message, fmt.Sprintf("%s: %s", "Subject", request.Body.Subject))
	message = append(message, fmt.Sprintf("%s: %s", "MIME-Version", `1.0`))
	if len(request.Body.Files) > 0 {
		msg := fmt.Sprintf("%s", request.Body.Msg)
		boundary := api.Boundary()
		message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `multipart/mixed;boundary=`+boundary))
		message = append(message, fmt.Sprintf("--%s", boundary))
		message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset="utf-8"`))
		message = append(message, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
		message = append(message, base64.StdEncoding.EncodeToString([]byte(msg)))
		message = append(message, fmt.Sprintf("--%s", boundary))
		for _, file := range request.Body.Files {
			filename := fmt.Sprintf("%s", file)
			content, err := ioutil.ReadFile(filename)
			filetype := http.DetectContentType(content)
			if err == nil {
				fmt.Println(filetype)
				message = append(message, fmt.Sprintf("%s: %s", "Content-Type", filetype))
				message = append(message, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
				message = append(message, fmt.Sprintf("%s: %s", "Content-Disposition", `attachment; filename=`+path.Base(filename)))
				message = append(message, base64.StdEncoding.EncodeToString(content))
				message = append(message, fmt.Sprintf("--%s", boundary))
			} else {
				return false
			}
		}
	} else {
		msg := fmt.Sprintf("%s", request.Body.Msg)
		message = append(message, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset="utf-8"`))
		message = append(message, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
		message = append(message, base64.StdEncoding.EncodeToString([]byte(msg)))
	}
	auth := smtp.PlainAuth("", config.MailUser, config.MailPass, config.MailHost)
	addr := config.MailHost + ":" + config.MailPort
	err := smtp.SendMail(addr, auth, request.Body.From.Address, []string{request.Body.To.Address}, []byte(strings.Join(message, "\r\n")))
	if err != nil {
		return false
	}
	return true
}
