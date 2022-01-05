//xc
// 2022-01-04
// Digimaker mail that supports outlook mail
// Usage: import this package in project.
package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"

	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/jordan-wright/email"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func SendMail(mail util.MailMessage) error {
	from := util.GetConfig("general", "send_from")
	hostPort := util.GetConfig("general", "mail_host")
	password := util.GetConfig("general", "mail_password")
	host, _, _ := net.SplitHostPort(hostPort)

	message := email.NewEmail()

	message.To = mail.To
	message.From = from
	message.Subject = mail.Subject
	message.HTML = []byte(mail.Body)
	for _, attachment := range mail.Attachments {
		message.AttachFile(util.VarFolder() + "/" + attachment)
	}

	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	tlsconfig := &tls.Config{
		ServerName: host,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		return err
	}

	auth := LoginAuth(from, password)

	if err = c.Auth(auth); err != nil {
		return err
	}

	err = message.Send(hostPort, auth)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	util.HandleSendMail(func(mail util.MailMessage) error {
		err := SendMail(mail)
		if err != nil {
			log.Error(fmt.Errorf("Error when sending email: %v", err), "")
		}
		return err
	})
}
