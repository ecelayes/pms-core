package ports

import (
	"context"
	"github.com/ecelayes/pms-core/internal/core/domain"
)

type AvailabilityRepository interface {
	GetAvailability(ctx context.Context, startDate, endDate string) ([]domain.InventoryItem, error)
}

type BookingRepository interface {
	CreateBookingAtomic(ctx context.Context, req domain.BookingRequest) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Booking, error)
  List(ctx context.Context, hotelID string) ([]domain.Booking, error)
  UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error
}
