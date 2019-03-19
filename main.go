package main

import (
	"fmt"
	"net/mail"

	mailer "github.com/OzqurYalcin/mailer/src"
)

func main() {
	config := mailer.Config{MailHost: "", MailPort: "", MailUser: "", MailPass: ""}
	api := &mailer.API{Config: config}
	api.SetHeaders(mail.Address{"From Name", "from@example.com"}, mail.Address{"To Name", "to@example.com"}, "Title", "Message")
	// api.AttachFile("file.pdf")
	send := api.Send()
	if send {
		fmt.Println("SENT!")
	} else {
		fmt.Println("ERROR")
	}
}
