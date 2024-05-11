package common

type UserRegReq struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Password   string `json:"password"`
}

type UseLoginReq struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password" validate:"required"`
}
