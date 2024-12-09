package domain

type Account struct {
	ID       string   `json:"-"`
	Currency Currency `json:"currency"`
}
