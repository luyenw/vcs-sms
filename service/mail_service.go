package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"
	"vcs-sms/config/mail"
)

type MailService struct {
}

func NewMailService() *MailService {
	return &MailService{}
}
func (service *MailService) SendEmail(to []string, content string) error {
	auth := mail.GetAuth()
	title := "VCS-SMS Report " + time.Now().Format("2006-01-02")
	log.Printf("%+v\n", auth)
	header := make(map[string]string)
	header["From"] = "luyend785@gmail.com"
	header["To"] = strings.Join(to, ",")
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(content))
	return smtp.SendMail("smtp.gmail.com:587", auth, "luyend785@gmail.com", to, []byte(message))
}
