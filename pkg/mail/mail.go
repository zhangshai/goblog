package mail

import (
	"goblog/pkg/types"
	"mime"

	"gopkg.in/gomail.v2"
)

func SendMail(Mailto string, Subject string, Body string) (err error) {

	mailConn := map[string]string{
		"user": "657762708@qq.com",
		"pass": "dftklnfouczkbcbe",
		"host": "smtp.qq.com",
		"port": "465",
	}
	port := types.StringToInt(mailConn["port"])
	mail := gomail.NewMessage()
	mail.SetHeader("From", mime.QEncoding.Encode("UTF-8", "Support")+"<"+mailConn["user"]+">")
	mail.SetHeader("To", Mailto)
	mail.SetHeader("Subject", Subject)
	mail.SetBody("text/html", Body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	return d.DialAndSend(mail)

}
