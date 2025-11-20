package domain

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusPending   BookingStatus = "pending"
)

type BookingRequest struct {
	HotelID   string        `json:"hotel_id"`
	Guest     GuestInfo     `json:"guest"`
	Items     []BookingItem `json:"items"`
	Currency  string        `json:"currency"`
}

type GuestInfo struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type BookingItem struct {
	RoomTypeID string `json:"room_type_id"`
	RatePlanID string `json:"rate_plan_id"`
	CheckIn    string `json:"check_in"`  // YYYY-MM-DD
	CheckOut   string `json:"check_out"` // YYYY-MM-DD
	Quantity   int    `json:"quantity"`
}

type Booking struct {
	ID          string        `json:"id"`
	Code        string        `json:"code"`
	Status      BookingStatus `json:"status"`
	GuestID     string        `json:"guest_id,omitempty"`
	GuestName   string        `json:"guest_name,omitempty"`
	TotalAmount float64       `json:"total_amount"`
	Currency    string        `json:"currency"`
	CreatedAt   time.Time     `json:"created_at"`
	Guest       *GuestInfo    `json:"guest_details,omitempty"`
	Items       []BookingItem `json:"items,omitempty"`
}
