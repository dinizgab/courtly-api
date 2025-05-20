package entity

import "time"

type Court struct {
	ID          string  `json:"id"`
	CompanyId   string  `json:"company_id"`
	Name        string  `json:"name"`
	IsActive    bool    `json:"is_active"`
	SportType   string  `json:"sport_type"`
	HourlyPrice float64 `json:"hourly_price"`
	Description string  `json:"description"`
	Capacity    int     `json:"capacity"`
    OpeningTime time.Time `json:"opening_time"`
    ClosingTime time.Time `json:"closing_time"`

	Bookings []Booking `json:"bookings"`
}
