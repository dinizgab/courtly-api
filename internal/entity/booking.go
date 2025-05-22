package entity

import (
	"errors"
	"math/rand/v2"
	"time"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

var (
	ErrInvalidVerificationCode = errors.New("invalid verification code")
	ErrBookingAlreadyConfirmed = errors.New("booking already confirmed")
	ErrInvalidCodeFormat       = errors.New("verification code must be 6 digits")
)

type Booking struct {
	ID               string        `json:"id"`
	CourtId          string        `json:"court_id"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
	CreatedAt        time.Time     `json:"created_at"`
	Status           BookingStatus `json:"status"`
	GuestName        string        `json:"guest_name"`
	GuestPhone       string        `json:"guest_phone"`
	GuestEmail       string        `json:"guest_email"`
	VerificationCode string        `json:"verification_code"`
	Court            *Court        `json:"court,omitempty"`
}

func GenerateVerificationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)

	for i := 0; i < 6; i++ {
		code[i] = charset[rand.IntN(len(charset))]
	}

	return string(code)
}
