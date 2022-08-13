# Smtp.go
Simple library for sending utf-8 e-mails via smtp with golang

# License
The MIT License (MIT). Please see License File for more information.

# Installation
```bash
go get github.com/ozgur-soft/smtp.go
```

# Usage
```go
package main

import (
	"fmt"
	"net/mail"

	smtp "github.com/ozgur-soft/smtp.go/src"
)

func main() {
	config := smtp.Config{MailHost: "", MailPort: "", MailUser: "", MailPass: ""}
	api := &smtp.API{Config: config}
	api.SetHeaders(mail.Address{"From Name", "from@example.com"}, mail.Address{"To Name", "to@example.com"}, "Title", "Message")
	send := api.Send()
	if send {
		fmt.Println("SENT!")
	} else {
		fmt.Println("ERROR")
	}
}
```
