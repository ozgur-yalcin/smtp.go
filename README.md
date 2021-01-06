# Mailer
Simple library for sending utf-8 e-mails and attachments via smtp with golang

# License
The MIT License (MIT). Please see License File for more information.

# Installation
```bash
go get github.com/ozgur-soft/mailer
```

# Usage
```go
package main

import (
	"fmt"
	"net/mail"

	mailer "github.com/ozgur-soft/mailer/src"
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
```
