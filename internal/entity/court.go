package entity

import "time"

type Court struct {
	ID        string `json:"id"`
	CompanyId string `json:"company_id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	SportType string `json:"sport_type"`
	// TODO - Change to int64 for cents
	HourlyPrice   float64       `json:"hourly_price"`
	Description   string        `json:"description"`
	Capacity      int           `json:"capacity"`
	BookingsToday int           `json:"bookings_today"`
	Bookings      []Booking     `json:"bookings"`
	Company       *Company      `json:"company,omitempty"`
	Photos        []CourtPhoto  `json:"photos,omitempty"`
	CourtSchedule []CourtSchedule `json:"court_schedule"`
}

type CourtPhoto struct {
	ID       string `json:"id"`
	CourtId  string `json:"court_id"`
	Path     string `json:"path"`
	Position int    `json:"position"`
	IsCover  bool   `json:"is_cover"`
}

type CourtSchedule struct {
	CourtId     string    `json:"court_id"`
    Weekday     int       `json:"weekday"`
	IsOpen      bool      `json:"is_open"`
	OpeningTime time.Time `json:"opening_time"`
	ClosingTime time.Time `json:"closing_time"`
}
