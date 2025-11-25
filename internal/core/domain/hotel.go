package domain

type Hotel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
	Currency  string `json:"currency"`
	OwnerID   string `json:"owner_id"`
}

type CreateHotelRequest struct {
	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`
	Currency  string `json:"currency"`
}
