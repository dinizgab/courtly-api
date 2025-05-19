package entity

type Court struct {
	ID          string  `json:"id"`
	CompanyId   string  `json:"company_id"`
	Name        string  `json:"name"`
	IsActive    bool    `json:"is_active"`
	SportType   string  `json:"sport_type"`
	HourlyPrice float64 `json:"hourly_price"`

	Bookings []Booking `json:"bookings"`
}
