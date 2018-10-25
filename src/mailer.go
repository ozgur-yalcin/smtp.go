package mailer

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/mail"
	"net/smtp"
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
	buffer := bytes.NewBuffer(nil)
	boudary := api.Boundary()
	breakline := "\r\n"
	if len(request.Body.Files) > 0 {
		header := []string{}
		header = append(header, fmt.Sprintf("%s: %s", "To", request.Body.To.String()))
		header = append(header, fmt.Sprintf("%s: %s", "From", request.Body.From.String()))
		header = append(header, fmt.Sprintf("%s: %s", "Subject", request.Body.Subject))
		header = append(header, fmt.Sprintf("%s: %s", "Mime-Version", "1.0"))
		header = append(header, fmt.Sprintf("%s: %s", "Content-Type", `multipart/mixed;boundary="`+boudary+`"`))
		header = append(header, fmt.Sprintf("%s", ""))
		header = append(header, fmt.Sprintf("--%s", boudary))
		header = append(header, fmt.Sprintf("%s", ""))
		buffer.WriteString(strings.Join(header, breakline))
		content := []string{}
		content = append(content, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset=UTF-8`))
		content = append(content, fmt.Sprintf("%s", ""))
		content = append(content, fmt.Sprintf("%s", request.Body.Msg))
		buffer.WriteString(strings.Join(content, breakline))
		for _, filename := range request.Body.Files {
			f := fmt.Sprintf("%s", filename)
			file, _ := ioutil.ReadFile(f)
			attachment := []string{}
			attachment = append(attachment, fmt.Sprintf("%s", ""))
			attachment = append(attachment, fmt.Sprintf("--%s", boudary))
			attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
			attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Disposition", `attachment`))
			attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Type", http.DetectContentType(file)+`;name="`+f+`"`))
			attachment = append(attachment, fmt.Sprintf("%s", ""))
			buffer.WriteString(strings.Join(attachment, breakline))
			b := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
			base64.StdEncoding.Encode(b, file)
			buffer.WriteString(breakline)
			for i, l := 0, len(b); i < l; i++ {
				buffer.WriteByte(b[i])
				if (i+1)%76 == 0 {
					buffer.WriteString(breakline)
				}
			}
		}
		buffer.WriteString(breakline)
		buffer.WriteString("--" + boudary + "--")
	} else {
		content := []string{}
		content = append(content, fmt.Sprintf("%s: %s", "To", request.Body.To.String()))
		content = append(content, fmt.Sprintf("%s: %s", "From", request.Body.From.String()))
		content = append(content, fmt.Sprintf("%s: %s", "Subject", request.Body.Subject))
		content = append(content, fmt.Sprintf("%s: %s", "Mime-Version", "1.0"))
		content = append(content, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset=UTF-8`))
		content = append(content, fmt.Sprintf("%s", ""))
		content = append(content, fmt.Sprintf("%s", request.Body.Msg))
		buffer.WriteString(strings.Join(content, breakline))
	}
	auth := smtp.PlainAuth("", config.MailUser, config.MailPass, config.MailHost)
	addr := config.MailHost + ":" + config.MailPort
	err := smtp.SendMail(addr, auth, request.Body.From.Address, []string{request.Body.To.Address}, buffer.Bytes())
	if err != nil {
		return false
	}
	return true
}
