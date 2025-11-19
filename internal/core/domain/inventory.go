package domain

type InventoryItem struct {
	RoomTypeID string  `json:"room_type_id"`
	Date       string  `json:"date"`
	Available  int     `json:"available"`
	Price      float64 `json:"price"`
}
