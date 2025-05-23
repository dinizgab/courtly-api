package notification

import (
	"context"

	"github.com/dinizgab/booking-mvp/internal/config"
	"github.com/dinizgab/booking-mvp/internal/entity"
	"gopkg.in/gomail.v2"
)

const emailSubject = "Confirmação da reserva"

type Sender interface {
	Send(ctx context.Context, subject entity.BookingConfirmationDTO) error
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

func (s *emailSender) Send(ctx context.Context, info entity.BookingConfirmationDTO) error {
	body, err := s.Renderer.Render(info)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.Config.Email)
	m.SetHeader("To", info.GuestEmail)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.Config.Host, s.Config.Port, s.Config.User, s.Config.Pass)

	return d.DialAndSend(m)
}
