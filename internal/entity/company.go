package entity

type Company struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	CNPJ 	 string `json:"cnpj"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Slug     string `json:"slug"`

	Courts []Court `json:"courts"`
}
