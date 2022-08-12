package mailer

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"path"
	"strings"
)

type Config struct {
	MailHost string
	MailPort string
	MailUser string
	MailPass string
}

type API struct {
	Buffer   *bytes.Buffer
	Boundary interface{}
	Config   Config
	Header   []string
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
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "To", api.Body.To.String()))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "From", api.Body.From.String()))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "Subject", api.Body.Subject))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "Mime-Version", "1.0"))
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", "Content-Type", `multipart/mixed;boundary="`+api.Boundary.(string)+`"`))
	api.Header = append(api.Header, fmt.Sprintf("%s", ""))
	api.Header = append(api.Header, fmt.Sprintf("--%s", api.Boundary.(string)))
	api.Header = append(api.Header, fmt.Sprintf("%s", ""))
	api.Buffer.WriteString(strings.Join(api.Header, "\r\n"))
	content := []string{}
	content = append(content, fmt.Sprintf("%s: %s", "Content-Type", `text/html;charset=UTF-8`))
	content = append(content, fmt.Sprintf("%s", ""))
	content = append(content, fmt.Sprintf("%s", api.Body.Message))
	api.Buffer.WriteString(strings.Join(content, "\r\n"))
	return
}

func (api *API) AddHeader(key, value string) {
	api.Header = append(api.Header, fmt.Sprintf("%s: %s", key, value))
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
			return api.AttachFileBytes(path.Base(filepath), file)
		} else {
			return false
		}
	} else {
		return false
	}
}

func (api *API) AttachFileBytes(name string, file []byte) bool {
	if api.Boundary != nil {
		attachment := []string{}
		attachment = append(attachment, fmt.Sprintf("%s", ""))
		attachment = append(attachment, fmt.Sprintf("--%s", api.Boundary.(string)))
		attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Transfer-Encoding", `base64`))
		attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Disposition", `attachment`))
		attachment = append(attachment, fmt.Sprintf("%s: %s", "Content-Type", http.DetectContentType(file)+`;name="`+name+`"`))
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
}

func (api *API) Send() bool {
	if api.Boundary != nil {
		api.Buffer.WriteString("\r\n")
		api.Buffer.WriteString("--" + api.Boundary.(string) + "--")
		addr := api.Config.MailHost + ":" + api.Config.MailPort
		host, _, _ := net.SplitHostPort(addr)
		auth := smtp.PlainAuth("", api.Config.MailUser, api.Config.MailPass, host)
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
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
	} else {
		return false
	}
}
