package main

import (
	"fmt"
	"mailer/config"
	"mailer/src"
	"net/mail"
)

func init() {
	config.MailHost = "" // Mail Host
	config.MailPort = "" // Mail Port
	config.MailUser = "" // Kullanıcı Adı
	config.MailPass = "" // Şifre
}

func main() {
	maildata := mailer.MailData{}
	maildata.Body.From = mail.Address{"Name", "mail@example.com"}
	maildata.Body.To = mail.Address{"Name", "mail@example.com"}
	maildata.Body.Subject = "Title"
	maildata.Body.Msg = "Message"
	send := mailer.Mail(maildata)
	if send {
		fmt.Println("SENT")
	} else {
		fmt.Println("ERROR")
	}
}
