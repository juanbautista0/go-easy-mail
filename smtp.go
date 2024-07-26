package go_easy_mail

import (
	"crypto/tls"
	"errors"
	"net/smtp"
	"regexp"
	"strconv"
)

type SMTPClientImpl struct {
	serverName string
	tlsConfig  *tls.Config
	Debug      bool
	Host       string
	Port       string
}
type SMTPCredentialClient struct {
	Host     string
	Port     string
	UserName string
	Password string
}

func newSMTPClient(host, port string, tlsConfig ...*tls.Config) *SMTPClientImpl {
	instance := &SMTPClientImpl{
		serverName: host + ":" + port,
		Debug:      false,
		Host:       host,
		Port:       port,
		tlsConfig: &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         host,
		},
	}

	if len(tlsConfig) > 0 {
		instance.tlsConfig = tlsConfig[0]
	}

	return instance
}

func (s *SMTPClientImpl) conn() (*tls.Conn, error) {
	var conn *tls.Conn
	var err error
	if conn, err = tls.Dial("tcp", s.serverName, s.tlsConfig); err != nil {
		return nil, err
	}

	return conn, nil
}

func (s *SMTPClientImpl) getClient() (*smtp.Client, error) {
	var conn *tls.Conn
	var err error
	var client *smtp.Client
	var portNumber int
	if !s.isDomain(s.Host) {
		return nil, errors.New(s.Host + " is NOT a valid SMTP host.")
	}

	if portNumber, err = strconv.Atoi(s.Port); err != nil || portNumber <= 0 {
		return nil, errors.New(s.Port + " is NOT a valid SMTP port.")
	}

	if conn, err = s.conn(); err != nil {
		return nil, err
	}

	if client, err = smtp.NewClient(conn, s.Host); err != nil {
		return nil, err
	}

	return client, nil
}

func (s *SMTPClientImpl) isDomain(domain string) bool {
	var domainRegex = regexp.MustCompile(`^(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,}$`)
	return domainRegex.MatchString(domain)
}
