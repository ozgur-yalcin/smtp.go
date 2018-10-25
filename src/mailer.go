package mailer

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
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
	Buffer   *bytes.Buffer
	Boundary interface{}
	Body     struct {
		To      mail.Address
		From    mail.Address
		Subject interface{}
		Message interface{}
	}
}

func (api *API) SetHeaders(from, to mail.Address, subject, message interface{}) {
	api.Buffer = bytes.NewBuffer(nil)
	api.SetBoundary()
	api.Body.From = from
	api.Body.To = to
	api.Body.Subject = subject
	api.Body.Message = message
	header := []string{}
	header = append(header, fmt.Sprintf("%s: %s", "To", api.Body.To.String()))
	header = append(header, fmt.Sprintf("%s: %s", "From", api.Body.From.String()))
	header = append(header, fmt.Sprintf("%s: %s", "Subject", api.Body.Subject))
	header = append(header, fmt.Sprintf("%s: %s", "Mime-Version", "1.0"))
	header = append(header, fmt.Sprintf("%s: %s", "Content-Type", `multipart/mixed;boundary="`+api.Boundary.(string)+`"`))
	header = append(header, fmt.Sprintf("%s", ""))
	header = append(header, fmt.Sprintf("--%s", api.Boundary.(string)))
	header = append(header, fmt.Sprintf("%s", ""))
	api.Buffer.WriteString(strings.Join(header, "\r\n"))
	content := []string{}
	content = append(content, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset=UTF-8`))
	content = append(content, fmt.Sprintf("%s", ""))
	content = append(content, fmt.Sprintf("%s", api.Body.Message))
	api.Buffer.WriteString(strings.Join(content, "\r\n"))
	return
}

func (api *API) SetBoundary() {
	var buf [30]byte
	io.ReadFull(rand.Reader, buf[:])
	api.Boundary = fmt.Sprintf("%x", buf[:])
	return
}

func (api *API) AttachFile(filepath string) bool {
	if api.Boundary != nil {
		file, err := ioutil.ReadFile(filepath)
		if err == nil {
			attachment := []string{}
			attachment = append(attachment, fmt.Sprintf("%s", ""))
			attachment = append(attachment, fmt.Sprintf("--%s", api.Boundary.(string)))
			attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
			attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Disposition", `attachment`))
			attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Type", http.DetectContentType(file)+`;name="`+path.Base(filepath)+`"`))
			attachment = append(attachment, fmt.Sprintf("%s", ""))
			api.Buffer.WriteString(strings.Join(attachment, "\r\n"))
			b := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
			base64.StdEncoding.Encode(b, file)
			api.Buffer.WriteString("\r\n")
			for i, l := 0, len(b); i < l; i++ {
				api.Buffer.WriteByte(b[i])
				if (i+1)%76 == 0 {
					api.Buffer.WriteString("\r\n")
				}
			}
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (api *API) Send() bool {
	if api.Boundary != nil {
		api.Buffer.WriteString("\r\n")
		api.Buffer.WriteString("--" + api.Boundary.(string) + "--")
		auth := smtp.PlainAuth("", config.MailUser, config.MailPass, config.MailHost)
		addr := config.MailHost + ":" + config.MailPort
		err := smtp.SendMail(addr, auth, api.Body.From.Address, []string{api.Body.To.Address}, api.Buffer.Bytes())
		if err != nil {
			return false
		}
		return true
	} else {
		return false
	}
}
