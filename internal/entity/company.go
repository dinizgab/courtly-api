package entity

type Company struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Slug    string `json:"slug"`

	Courts []Court `json:"courts"`
}
