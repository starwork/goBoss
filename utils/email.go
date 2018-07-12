package utils

import (
	"gopkg.in/gomail.v2"
	cf "goBoss/config"
	"log"
)

type Mail struct {
	Content  string
	Subject  string
	Attach   string
}

func (m *Mail) Send() {
	handle := gomail.NewMessage()
	handle.SetHeader("From", cf.Config.Sender)
	handle.SetHeader("To", cf.Config.Receiver)
	handle.SetHeader("Subject", "[Auto]:来自goBoss--"+m.Subject)
	handle.SetBody("text/html", m.Content)
	if m.Attach != "" {
		handle.Attach(m.Attach)
	}
	s := gomail.NewDialer(cf.Config.MailServer, 25, cf.Config.Sender, cf.Config.SenderPwd)
	if err := s.DialAndSend(handle); err != nil {
		log.Println("发送邮件失败!Error: ", err.Error())
	}
}
