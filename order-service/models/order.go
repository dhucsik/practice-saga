package models

type Order struct {
	ID     int    `json:"id"`
	ItemID int    `json:"item_id"`
	Status string `json:"status"`
	// suppose we have another fields

	IsPaid           bool `json:"-"`
	NotificationSent bool `json:"-"`
}
