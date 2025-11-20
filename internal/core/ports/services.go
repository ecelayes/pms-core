package ports

import (
	"context"
	"github.com/ecelayes/pms-core/internal/core/domain"
)

type BookingService interface {
	Create(ctx context.Context, req domain.BookingRequest) (string, error)
}
