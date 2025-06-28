package notification

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/config"
	"gopkg.in/gomail.v2"
)

type Sender interface {
	Send(ctx context.Context, tplName string, subject string, data any, to ...string) error
}

type emailSender struct {
	Config   *config.SMTPConfig
	Renderer Renderer
}

func NewEmailSender(renderer Renderer, config *config.SMTPConfig) Sender {
	return &emailSender{
		Renderer: renderer,
		Config:   config,
	}
}

func (s *emailSender) Send(
	ctx context.Context,
	tplName string,
	subject string,
	data any,
    to ...string,
) error {
	body, err := s.Renderer.Render(tplName, data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.Config.Email)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.Config.Host, s.Config.Port, s.Config.User, s.Config.Pass)

	return d.DialAndSend(m)
}
