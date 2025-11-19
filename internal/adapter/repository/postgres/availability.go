package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ecelayes/pms-core/internal/adapter/repository/postgres/db"

	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
)

type AvailabilityRepo struct {
	queries *db.Queries
}

func NewAvailabilityRepo(pool *pgxpool.Pool) ports.AvailabilityRepository {
	return &AvailabilityRepo{
		queries: db.New(pool),
	}
}

func (r *AvailabilityRepo) GetAvailability(ctx context.Context, start, end string) ([]domain.InventoryItem, error) {
	startDate, err := parseDate(start)
	if err != nil {
		return nil, err
	}

	endDate, err := parseDate(end)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetInventoryByDateRange(ctx, db.GetInventoryByDateRangeParams{
		Date:   startDate,
		Date_2: endDate,
	})
	if err != nil {
		return nil, err
	}

	items := make([]domain.InventoryItem, len(rows))

	for i, row := range rows {
		booked := int32(0)
		if row.BookedCount.Valid {
			booked = row.BookedCount.Int32
		}

		available := int(row.TotalInventory - booked)
		if available < 0 {
			available = 0
		}

		dateStr := row.Date.Time.Format("2006-01-02")
		price, _ := row.Price.Float64Value()

		items[i] = domain.InventoryItem{
			RoomTypeID: row.RoomTypeID.String(),
			Date:       dateStr,
			Available:  available,
			Price:      price.Float64,
		}
	}

	return items, nil
}

func parseDate(d string) (pgtype.Date, error) {
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		return pgtype.Date{}, err
	}
	return pgtype.Date{Time: t, Valid: true}, nil
}
