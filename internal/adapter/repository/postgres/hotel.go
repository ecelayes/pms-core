package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ecelayes/pms-core/internal/adapter/repository/postgres/db"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
)

type HotelRepo struct {
	pool *pgxpool.Pool
}

func NewHotelRepo(pool *pgxpool.Pool) ports.HotelRepository {
	return &HotelRepo{pool: pool}
}

func (r *HotelRepo) CreateHotel(ctx context.Context, ownerID string, req domain.CreateHotelRequest) (string, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return "", fmt.Errorf("invalid owner uuid: %w", err)
	}

	hotelID, err := db.New(r.pool).CreateHotel(ctx, db.CreateHotelParams{
		Name:      req.Name,
		Subdomain: req.Subdomain,
		Currency:  pgtype.Text{String: req.Currency, Valid: true},
		OwnerID:   pgtype.UUID{Bytes: ownerUUID, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("error creating hotel: %w", err)
	}
	return hotelID.String(), nil
}

func (r *HotelRepo) ListByOwner(ctx context.Context, ownerID string) ([]domain.Hotel, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner uuid: %w", err)
	}

	rows, err := db.New(r.pool).ListHotelsByOwner(ctx, pgtype.UUID{Bytes: ownerUUID, Valid: true})
	if err != nil {
		return nil, err
	}

	var hotels []domain.Hotel
	for _, row := range rows {
		hotels = append(hotels, domain.Hotel{
			ID:        row.ID.String(),
			Name:      row.Name,
			Subdomain: row.Subdomain,
			Currency:  row.Currency.String,
			OwnerID:   ownerID,
		})
	}
	return hotels, nil
}
