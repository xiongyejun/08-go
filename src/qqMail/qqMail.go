// https://studygolang.com/topics/2877
// user:登陆邮箱账号
// password:不是qq邮箱密码,需要登陆你的qq邮箱，在设置，账号，启用IMAP/SMTP服务，会发送一段身份验证符号给你，用这个登陆
// host:smtp.qq.com:587
// to:加入多个邮箱,已逗号隔开,相当于群发。
// subject:发送的主题
// body:发送的内容
// mailtyoe: 发送的内容是文本还是html

package main

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"

	"github.com/scorredoira/email"
)

type myMmail struct {
	user     string
	password string
	host     string
	to       []string
	subject  string
	body     string

	attach []string
}

func (me *myMmail) mailInit() {
	me.user = "648555205@qq.com"
	bpassword, _ := ioutil.ReadFile("E:/648555205mailKey.txt")
	me.password = string(bpassword)
	me.host = "smtp.qq.com:587"
	me.to = append(me.to, "648555205@qq.com")
	me.to = append(me.to, "244114746@qq.com")

	me.subject = "Test send email by golang"
	me.body = "Test send email by golang"

	me.attach = append(me.attach, "qqMail.go")
	me.attach = append(me.attach, "qqMail1.txt")

}

//func SendMail(user, password, host, to, subject, body, mailtype string) error {
//	hp := strings.Split(host, ":")
//	auth := smtp.PlainAuth(password, user, password, hp[0])
//	var content_type string
//	if mailtype == "html" {
//		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
//	} else {
//		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
//	}
//	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
//	send_to := strings.Split(to, ";")

//	err := smtp.SendMail(host, auth, user, send_to, msg)
//	return err
//}

func main() {
	m := new(myMmail)
	m.mailInit()

	msg := email.NewMessage(m.subject, m.body)
	msg.From.Address = m.user
	msg.To = m.to

	for _, v := range m.attach {
		if err := msg.Attach(v); err != nil {
			fmt.Println(1, err)
			return
		}
	}

	hp := strings.Split(m.host, ":")
	auth := smtp.PlainAuth(m.password, m.user, m.password, hp[0])

	if err := email.Send(m.host, auth, msg); err != nil {
		fmt.Println(2, err)
		return
	}
}
