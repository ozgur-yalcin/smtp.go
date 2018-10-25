# Mailer
Simple library for sending utf-8 e-mail via smtp with golang

# License
The MIT License (MIT). Please see License File for more information.

# Installation
```bash
go get github.com/OzqurYalcin/mailer
```

# Usage
```go
package main

import (
	"fmt"
	"net/mail"

	"github.com/OzqurYalcin/mailer/config"
	"github.com/OzqurYalcin/mailer/src"
)

func init() {
	config.MailHost = "" // Mail Host
	config.MailPort = "" // Mail Port
	config.MailUser = "" // Mail Username
	config.MailPass = "" // Mail Password
}

func main() {
	api := new(mailer.API)
	api.Lock()
	defer api.Unlock()
	api.SetHeaders(mail.Address{"From", "from@example.com"}, mail.Address{"To", "to@example.com"}, "Title", "Message")
	// api.AttachFile("file.pdf")
	send := api.Send()
	if send {
		fmt.Println("SENT!")
	} else {
		fmt.Println("ERROR")
	}
}
```
