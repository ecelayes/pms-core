package ports

import (
	"context"
	"github.com/ecelayes/pms-core/internal/core/domain"
)

type BookingService interface {
	Create(ctx context.Context, req domain.BookingRequest) (string, error)
	Get(ctx context.Context, id string) (*domain.Booking, error)
  List(ctx context.Context, hotelID string) ([]domain.Booking, error)
  Cancel(ctx context.Context, id string) error
}
