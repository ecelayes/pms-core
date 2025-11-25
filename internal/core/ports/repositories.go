package ports

import (
	"context"
	"github.com/ecelayes/pms-core/internal/core/domain"
)

type AuthRepository interface {
	RegisterUser(ctx context.Context, email, passwordHash string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type HotelRepository interface {
	CreateHotel(ctx context.Context, ownerID string, req domain.CreateHotelRequest) (string, error)
	ListByOwner(ctx context.Context, ownerID string) ([]domain.Hotel, error)
}

type AvailabilityRepository interface {
	GetAvailability(ctx context.Context, startDate, endDate string) ([]domain.InventoryItem, error)
}

type BookingRepository interface {
	CreateBookingAtomic(ctx context.Context, req domain.BookingRequest) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
  List(ctx context.Context, hotelID string) ([]domain.Booking, error)
  UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error
	CancelBookingAtomic(ctx context.Context, id string) error
}
