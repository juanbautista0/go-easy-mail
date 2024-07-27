package go_easy_mail

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

type Mail struct {
	Sender      string
	SenderName  string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
	IsHTML      bool
}

type GoEasyEmail struct {
	host       string
	port       string
	userName   string
	password   string
	encryption bool
}

func NewGoEasyEmail(host, port, user, password string, encryption bool) *GoEasyEmail {
	return &GoEasyEmail{
		host:       host,
		port:       port,
		userName:   user,
		password:   password,
		encryption: encryption,
	}
}

func (in *GoEasyEmail) Send(mail *Mail) (*Mail, error) {
	var client *smtp.Client
	var err error
	var w io.WriteCloser
	var receivers []string
	var messageBody string = ""

	if _, err = in.validateMail(mail); err != nil {
		return mail, err
	}

	messageBody = in.buildMessage(mail)

	app := newSMTPClient(in.host, in.port)
	app.tlsConfig = &tls.Config{
		InsecureSkipVerify: in.encryption,
		ServerName:         in.host,
	}

	if client, err = app.getClient(); err != nil {
		return mail, err
	}

	//	use auth
	auth := smtp.PlainAuth("", in.userName, in.password, in.host)
	if err = client.Auth(auth); err != nil {
		return mail, err
	}

	//	add all from and to
	if err = client.Mail(mail.Sender); err != nil {
		return mail, err
	}
	receivers = append(mail.To, mail.Cc...)
	receivers = append(receivers, mail.Bcc...)

	for _, k := range receivers {

		if err = client.Rcpt(k); err != nil {
			log.Println(err.Error())
		}
	}

	// Data
	if w, err = client.Data(); err != nil {
		return mail, err
	}

	if _, err = w.Write([]byte(messageBody)); err != nil {
		return mail, err
	}

	err = w.Close()
	if err != nil {
		log.Println(err.Error())
	}

	client.Quit()

	//	Mail sent successfully
	return mail, nil
}

func (in *GoEasyEmail) buildMessage(mail *Mail) string {
	var message strings.Builder

	// MIME Header
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: multipart/mixed; boundary=\"myboundary\"\r\n")

	if mail.SenderName != "" {
		message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", mail.SenderName, mail.Sender))
	} else {
		message.WriteString(fmt.Sprintf("From: %s\r\n", mail.Sender))
	}

	if len(mail.To) > 0 {
		message.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	}
	if len(mail.Cc) > 0 {
		message.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";")))
	}
	subject := base64.StdEncoding.EncodeToString([]byte(mail.Subject))
	message.WriteString(fmt.Sprintf("Subject: =?utf-8?B?%s?=\r\n", subject))
	message.WriteString("\r\n")

	// Message body
	message.WriteString("--myboundary\r\n")
	if mail.IsHTML {
		message.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	} else {
		message.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	}
	message.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	message.WriteString("\r\n")
	message.WriteString(mail.Body)
	message.WriteString("\r\n")

	// Attachments
	for filename, data := range mail.Attachments {
		message.WriteString("--myboundary\r\n")
		contentType := in.getContentType(filename)
		message.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", contentType, filename))
		message.WriteString("Content-Transfer-Encoding: base64\r\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", filename))
		message.WriteString("\r\n")
		message.WriteString(base64.StdEncoding.EncodeToString(data))
		message.WriteString("\r\n")
	}

	message.WriteString("--myboundary--")

	return message.String()
}

func (in *GoEasyEmail) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		return "text/plain"
	case ".html":
		return "text/html"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

func (in *GoEasyEmail) validateMail(mail *Mail) (bool, error) {
	errorEmailAddres := "email address is inavlid"

	if mail.Body == "" {
		return false, errors.New("must indicate a body")
	}

	if mail.Subject == "" {
		return false, errors.New("must indicate a subject")
	}

	if !in.IsEmail(mail.Sender) {
		return false, errors.New(errorEmailAddres)
	}

	if len(mail.To) <= 0 {
		return false, errors.New("you must indicate at least one mail recipient")
	}

	for _, e := range mail.To {
		if !in.IsEmail(e) {
			return false, errors.New(errorEmailAddres)
		}
	}

	for _, e := range mail.Cc {
		if !in.IsEmail(e) {
			return false, errors.New(errorEmailAddres)
		}
	}

	for _, e := range mail.Bcc {
		if !in.IsEmail(e) {
			return false, errors.New(errorEmailAddres)
		}
	}

	return true, nil
}

func (in *GoEasyEmail) IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (in *GoEasyEmail) ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}
	return data
}
