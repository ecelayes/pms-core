package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ecelayes/pms-core/internal/adapter/repository/postgres/db"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) ports.AuthRepository {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) RegisterUser(ctx context.Context, email, passHash string) (string, error) {
	userID, err := db.New(r.pool).CreateUser(ctx, db.CreateUserParams{
		HotelID:      pgtype.UUID{Valid: false},
		Email:        email,
		PasswordHash: passHash,
		Role:         pgtype.Text{String: "owner", Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("Error creating user: %w", err)
	}
	return userID.String(), nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := db.New(r.pool).GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var hotelID *string
	if u.HotelID.Valid {
		s := uuid.UUID(u.HotelID.Bytes).String()
		hotelID = &s
	}

	return &domain.User{
		ID:           u.ID.String(),
		HotelID:      hotelID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Role:         u.Role.String,
	}, nil
}
