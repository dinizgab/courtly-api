package entity

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
	ErrBookingNotFound         = errors.New("booking not found")
)

type BookingFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
}

type BookingConfirmationInfo struct {
	ID               string `json:"id"`
	GuestName        string `json:"guest_name"`
	GuestPhone       string `json:"guest_phone"`
	GuestEmail       string `json:"guest_email"`
	CourtName        string `json:"court_name"`
	CourtAddress     string `json:"court_address"`
	BookingDate      string `json:"booking_date"`
	BookingInterval  string `json:"booking_interval"`
	TotalPrice       string `json:"total_price"`
	VerificationCode string `json:"verification_code"`
	CancelToken      string `json:"cancel_token"`
}

type Booking struct {
	ID                       string        `json:"id"`
	CourtId                  string        `json:"court_id"`
	StartTime                time.Time     `json:"start_time"`
	EndTime                  time.Time     `json:"end_time"`
	CreatedAt                time.Time     `json:"created_at"`
	Status                   BookingStatus `json:"status"`
	GuestName                string        `json:"guest_name"`
	GuestPhone               string        `json:"guest_phone"`
	GuestEmail               string        `json:"guest_email"`
	VerificationCode         string        `json:"verification_code"`
	TotalPrice               int64         `json:"total_price"`
	CancelTokenHash          string        `json:"cancel_token_hash"`
	CancelTokenHashExpiresAt time.Time     `json:"cancel_token_hash_expires_at"`
	Court                    *Court        `json:"court,omitempty"`
}

func (b Booking) DurationInHours() float64 {
	if b.EndTime.IsZero() || b.StartTime.IsZero() {
		return 0
	}
	duration := b.EndTime.Sub(b.StartTime).Hours()
	if duration < 0 {
		return 0
	}
	return duration
}

func GenerateVerificationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)

	for i := 0; i < 6; i++ {
		code[i] = charset[rand.IntN(len(charset))]
	}

	return string(code)
}

func GenerateCancelToken() (string, error) {
	b := make([]byte, 32)
	if _, err := crand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func HashCancelToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
