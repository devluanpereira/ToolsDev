package models

type Pagamento struct {
	ID         int
	UserID     int
	Email      string
	Quantidade int
	Status     string
	Link       string
}
