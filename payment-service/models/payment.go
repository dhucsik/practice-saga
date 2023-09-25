package models

type Payment struct {
	ID      int    `json:"id" db:"id"`
	OrderID int    `json:"order_id" db:"order_id"`
	Status  string `json:"status" db:"status"`
}
