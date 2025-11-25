package domain

type User struct {
	ID           string `json:"id"`
	HotelID      *string `json:"hotel_id,omitempty"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	// AdminName string `json:"admin_name"` // Opcional
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
