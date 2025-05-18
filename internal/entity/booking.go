package entity

import "time"

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
)

type Booking struct {
	ID               string        `json:"id"`
	CourtId          string        `json:"court_id"`
	Court            Court         `json:"court"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
	CreatedAt        time.Time     `json:"created_at"`
	Status           BookingStatus `json:"status"`
	GuestName        string        `json:"guest_name"`
	GuestPhone       string        `json:"guest_phone"`
	GuestEmail       string        `json:"guest_email"`
	VerificationCode string        `json:"verification_code"`
}
