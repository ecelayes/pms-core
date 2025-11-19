package ports

import (
	"context"
	"github.com/ecelayes/pms-core/internal/core/domain"
)

type AvailabilityRepository interface {
	GetAvailability(ctx context.Context, startDate, endDate string) ([]domain.InventoryItem, error)
}
