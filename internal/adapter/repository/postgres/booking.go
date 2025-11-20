package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ecelayes/pms-core/internal/adapter/repository/postgres/db"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
)

type BookingRepo struct {
	pool *pgxpool.Pool
}

func NewBookingRepo(pool *pgxpool.Pool) ports.BookingRepository {
	return &BookingRepo{pool: pool}
}

func (r *BookingRepo) CreateBookingAtomic(ctx context.Context, req domain.BookingRequest) (string, error) {
	hotelUUID, err := uuid.Parse(req.HotelID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", domain.ErrInvalidHotelID, err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	qtx := db.New(r.pool).WithTx(tx)

	guestID, err := qtx.CreateGuest(ctx, db.CreateGuestParams{
		HotelID:  pgtype.UUID{Bytes: hotelUUID, Valid: true},
		Email:    pgtype.Text{String: req.Guest.Email, Valid: true},
		FullName: req.Guest.FullName,
		Phone:    pgtype.Text{String: req.Guest.Phone, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("Error creating guest: %w", err)
	}

	var totalAmount float64
	bookingCode := generateResCode()

	if len(req.Items) == 0 {
		return "", errors.New("There are no items in the reservation.")
	}
	item := req.Items[0]
	
	start, _ := parseDate(item.CheckIn)
	end, _ := parseDate(item.CheckOut)
	
	roomTypeUUID, err := uuid.Parse(item.RoomTypeID)
	if err != nil {
		return "", fmt.Errorf("Invalid room_type_id: %w", err)
	}

	inventoryRows, err := qtx.GetInventoryByDateRange(ctx, db.GetInventoryByDateRangeParams{
		Date:   start,
		Date_2: end,
	})
	if err != nil {
		return "", fmt.Errorf("Error consulting inventory: %w", err)
	}

	for _, inv := range inventoryRows {
		if inv.RoomTypeID.Bytes != roomTypeUUID {
			continue
		}

		_, err := qtx.UpdateInventoryCount(ctx, db.UpdateInventoryCountParams{
			ID:      inv.ID,
			Version: inv.Version,
		})

		if err != nil {
			return "", fmt.Errorf("%w: date %s", domain.ErrNoAvailability, inv.Date.Time.Format("2006-01-02"))
		}

		p, _ := inv.Price.Float64Value()
		totalAmount += p.Float64
	}

	bookingID, err := qtx.CreateBooking(ctx, db.CreateBookingParams{
		HotelID:      pgtype.UUID{Bytes: hotelUUID, Valid: true},
		GuestID:      guestID,
		Code:         bookingCode,
		Status:       pgtype.Text{String: "confirmed", Valid: true},
		RoomAmount:   numericFromFloat(totalAmount),
		ExtrasAmount: numericFromFloat(0),
		TotalAmount:  numericFromFloat(totalAmount),
		Currency:     pgtype.Text{String: req.Currency, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("Error creating booking: %w", err)
	}

	err = qtx.CreateBookingRoom(ctx, db.CreateBookingRoomParams{
		BookingID:     bookingID.ID,
		RoomTypeID:    pgtype.UUID{Bytes: roomTypeUUID, Valid: true},
		CheckIn:       start,
		CheckOut:      end,
		PricePerNight: numericFromFloat(100.00),
		GuestNames:    []string{req.Guest.FullName},
	})
	if err != nil {
		return "", fmt.Errorf("Error creating booking detail: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", err
	}

	return bookingCode, nil
}

func generateResCode() string {
	return fmt.Sprintf("RES-%d", time.Now().UnixNano()%100000)
}

func numericFromFloat(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	n.Scan(fmt.Sprintf("%f", f))
	return n
}
