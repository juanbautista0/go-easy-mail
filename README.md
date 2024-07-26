# go-easy-mail

is a Go module designed to simplify sending emails via SMTP with support for SSL and TLS. This module allows developers to easily send emails, including the option to add attachments, all with a simple and easy to use interface.

### Features

- Emailing: Send emails quickly and easily.
- SSL/TLS support: Ensure the security of your emails through secure connections.
- Attachments: Add attachments to your e-mails without hassle.
- SMTP authentication: Supports the most common SMTP authentication methods.
- Simple interface: Designed to be intuitive and easy to use, even for those new to Go.

## Installation

```
go get github.com/juanbautista0/go-easy-mail

```


### Basic use

```go
package main

import go_easy_mail "github.com/juanbautista0/go-easy-mail"

func main() {
	host := "mail.yourmailhost.com"
	port := "465"
	user := "dev@yourmailhost.com"
	pass := "*********"
	encryption := true

	easyMail := go_easy_mail.NewGoEasyEmail(host, port, user, pass, encryption)

	mail := &go_easy_mail.Mail{}
	mail.Sender = user
	mail.To = []string{"recipient@easymail.com"}
	mail.Subject = "Test Email"
	mail.Body = "This is a test email using go-easy-mail."

	if _, err := easyMail.Send(mail); err == nil {
		fmt.Println("Mail sent successfully")
	}

}
```

### Attachments

```go
package main

import go_easy_mail "github.com/juanbautista0/go-easy-mail"

func main() {
	host := "mail.yourmailhost.com"
	port := "465"
	user := "dev@yourmailhost.com"
	pass := "*********"
	encryption := true

	easyMail := go_easy_mail.NewGoEasyEmail(host, port, user, pass, encryption)

	mail := &go_easy_mail.Mail{}
	mail.Sender = user
	mail.To = []string{"recipient@easymail.com"}
	mail.Subject = "Test Email"
	mail.Body = "This is a test email using go-easy-mail."

    //  set attachments
	mail.Attachments = map[string][]byte{
		"/tmp/MY_CV.pdf": easyMail.ReadFile("/tmp/MY_CV.pdf"),
	}

	if _, err := easyMail.Send(mail); err == nil {
		fmt.Println("Mail sent successfully")
	}

}
```