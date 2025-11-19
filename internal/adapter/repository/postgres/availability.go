package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
)

type AvailabilityRepo struct {
	db *pgxpool.Pool
}

func NewAvailabilityRepo(db *pgxpool.Pool) ports.AvailabilityRepository {
	return &AvailabilityRepo{db: db}
}

func (r *AvailabilityRepo) GetAvailability(ctx context.Context, start, end string) ([]domain.InventoryItem, error) {
	// TODO: We will put the actual SQL with sqlc here later.
	// For now, we return dummy data to test that the API responds.
	return []domain.InventoryItem{
		{
			RoomTypeID: "uuid-dummy",
			Date:       start,
			Available:  5,
			Price:      100.50,
		},
	}, nil
}
