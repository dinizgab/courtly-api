package entity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Company struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CNPJ     string `json:"cnpj"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Slug     string `json:"slug"`

    PixKey string `json:"pix_key"`

	Courts []Court `json:"courts"`
}

type CompanyDashboard struct {
	TotalBookings   int     `json:"total_bookings"`
	TotalEarnings   int64 `json:"total_earnings"`
	TotalClients    int     `json:"total_clients"`
	TotalBookedHours float64     `json:"total_booked_hours"`
}
