# mailer
Simple library for sending utf-8 e-mail via smtp with golang

# Security
If you discover any security related issues, please email ozguryalcin@outlook.com instead of using the issue tracker.

# License
The MIT License (MIT). Please see License File for more information.


```go
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
	api := new(mailer.API)
	api.Lock()
	defer api.Unlock()
	request := mailer.Request{}
	request.Body.From = mail.Address{"Name", "mail@example.com"}
	request.Body.To = mail.Address{"Name", "mail@example.com"}
	request.Body.Subject = "Title"
	request.Body.Msg = "Message"
	send := api.Mail(request)
	if send {
		fmt.Println("e-posta iletildi")
	} else {
		fmt.Println("hata oluştu")
	}
}
```
