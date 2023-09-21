package models

type Order struct {
	ID     int    `json:"id"`
	ItemID int    `json:"item_id"`
	Status string `json:"status"`
}
