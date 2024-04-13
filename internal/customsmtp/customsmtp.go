package customsmtp

import (
	"net/smtp"
	"os"
	"redGlow/internal/config"
	"strconv"
	"strings"
)

const HTMLTemplateReplace string = "--$$--"

type CustomSMTP struct {
	auth smtp.Auth
	cfg *config.Config
}

func NewSMTPService(cfg *config.Config) *CustomSMTP{
	return &CustomSMTP{
		auth: smtp.PlainAuth("",cfg.EmailSettings.HostUser,cfg.EmailSettings.HostPassword,cfg.EmailSettings.Host),
		cfg: cfg,
	}
}

func (cs *CustomSMTP) PrepareMessage(to_replace []string) ([]byte, error){
	HTMLBytes, err := os.ReadFile("./template/email_template.html")
	if err != nil {
		return nil,err
	}
	HTMLString := string(HTMLBytes)
	for _, item  := range to_replace {
		HTMLString = strings.Replace(HTMLString, HTMLTemplateReplace, item, 1)
	}
	return []byte(HTMLString), nil
}

func (cs *CustomSMTP) SendMail(send_to []string, message []byte) error {
	err := smtp.SendMail(
		cs.cfg.EmailSettings.Host+":"+strconv.Itoa(cs.cfg.EmailSettings.Port),
		cs.auth,
		cs.cfg.EmailSettings.HostUser,
		send_to,
		message,
	)
	if err != nil {
		return err
	}
	return nil
}