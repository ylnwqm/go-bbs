package utils

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

type Email struct {
	From     string
	To       string
	Header   map[string]string
	Body     string
	Template string
	Data     interface{}
	//CC		[]string
}

type Reg struct {
	Tag  string
	Code string
	Date string
}

func SendEmail(e Email) error {

	t, err := template.ParseFiles(e.Template)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	buffer := new(bytes.Buffer)
	t.Execute(buffer, e.Data)

	m := gomail.NewMessage()
	//m.SetHeader("From", e.From)
	// 增加发件人别名
	m.SetHeader("From",  m.FormatAddress(e.From, "霓虹灯下社区"))
	m.SetHeader("To", e.To)

	//m.SetAddressHeader("Cc", e.CC)
	m.SetHeader("Subject", e.Header["Subject"])
	m.SetBody("text/html", buffer.String())
	//m.Attach("/home/Alex/lolcat.jpg")
	// dbmlthxnihnrddji
	// itnkmmfmopneciha
	d := gomail.NewDialer("smtp.qq.com", 465, "1920853199@qq.com", "nczgumvhgqsodacb")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
