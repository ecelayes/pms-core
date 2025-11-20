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
	if err != nil { return "", err }
	defer tx.Rollback(ctx)

	qtx := db.New(r.pool).WithTx(tx)

	guestID, err := qtx.CreateGuest(ctx, db.CreateGuestParams{
		HotelID:  pgtype.UUID{Bytes: hotelUUID, Valid: true},
		Email:    pgtype.Text{String: req.Guest.Email, Valid: true},
		FullName: req.Guest.FullName,
		Phone:    pgtype.Text{String: req.Guest.Phone, Valid: true},
	})
	if err != nil { return "", fmt.Errorf("Error creating host: %w", err) }

	var totalAmount float64
	bookingCode := generateResCode()

	if len(req.Items) == 0 { return "", errors.New("There are no items.") }
	item := req.Items[0]
	
	start, _ := parseDate(item.CheckIn)
	end, _ := parseDate(item.CheckOut)
	roomTypeUUID, err := uuid.Parse(item.RoomTypeID)
	if err != nil { return "", fmt.Errorf("Invalid room_type_id: %w", err) }

	if item.Quantity < 1 { item.Quantity = 1 }

	inventoryRows, err := qtx.GetInventoryByDateRange(ctx, db.GetInventoryByDateRangeParams{
		Date:   start,
		Date_2: end,
	})
	if err != nil { return "", fmt.Errorf("Error checking inventory: %w", err) }

	for _, inv := range inventoryRows {
		if inv.RoomTypeID.Bytes != roomTypeUUID { continue }

		_, err := qtx.UpdateInventoryCount(ctx, db.UpdateInventoryCountParams{
			ID:       inv.ID,
			Version:  inv.Version,
			Quantity: int32(item.Quantity),
		})

		if err != nil {
			return "", fmt.Errorf("%w: There are not enough rooms available for that date %s", domain.ErrNoAvailability, inv.Date.Time.Format("2006-01-02"))
		}

		p, _ := inv.Price.Float64Value()
		totalAmount += p.Float64 * float64(item.Quantity)
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
	if err != nil { return "", fmt.Errorf("Error creating booking: %w", err) }

	for i := 0; i < item.Quantity; i++ {
		err = qtx.CreateBookingRoom(ctx, db.CreateBookingRoomParams{
			BookingID:     bookingID.ID,
			RoomTypeID:    pgtype.UUID{Bytes: roomTypeUUID, Valid: true},
			CheckIn:       start,
			CheckOut:      end,
			PricePerNight: numericFromFloat(100.00),
			GuestNames:    []string{req.Guest.FullName},
		})
		if err != nil { return "", fmt.Errorf("Error creating booking detail %d: %w", i+1, err) }
	}

	if err := tx.Commit(ctx); err != nil { return "", err }

	return bookingCode, nil
}

func (r *BookingRepo) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	bookingUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("Invalid uuid: %w", err)
	}

	row, err := db.New(r.pool).GetBookingByID(ctx, pgtype.UUID{Bytes: bookingUUID, Valid: true})
	if err != nil {
		return nil, err
	}

	amountVal, _ := row.TotalAmount.Float64Value()
	amount := amountVal.Float64

	booking := &domain.Booking{
		ID:          row.ID.String(),
		Code:        row.Code,
		Status:      domain.BookingStatus(row.Status.String),
		TotalAmount: amount,
		Currency:    row.Currency.String,
		CreatedAt:   row.CreatedAt.Time,
		Guest: &domain.GuestInfo{
			FullName: row.GuestName, 
			Email:    row.GuestEmail.String,
			Phone:    row.GuestPhone.String,
		},
	}

	rooms, err := db.New(r.pool).GetBookingRooms(ctx, pgtype.UUID{Bytes: bookingUUID, Valid: true})
	if err == nil {
		for _, room := range rooms {
			booking.Items = append(booking.Items, domain.BookingItem{
				RoomTypeID: room.RoomTypeID.String(),
				CheckIn:    room.CheckIn.Time.Format("2006-01-02"),
				CheckOut:   room.CheckOut.Time.Format("2006-01-02"),
				Quantity:   1,
			})
		}
	}

	return booking, nil
}

func (r *BookingRepo) List(ctx context.Context, hotelID string) ([]domain.Booking, error) {
	hID, _ := uuid.Parse(hotelID)
	
	rows, err := db.New(r.pool).ListBookings(ctx, db.ListBookingsParams{
		HotelID: pgtype.UUID{Bytes: hID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	var list []domain.Booking
	for _, row := range rows {
		amountVal, _ := row.TotalAmount.Float64Value()
		amount := amountVal.Float64

		list = append(list, domain.Booking{
			ID:          row.ID.String(),
			Code:        row.Code,
			Status:      domain.BookingStatus(row.Status.String),
			GuestName:   row.GuestName, 
			TotalAmount: amount,
			CreatedAt:   row.CreatedAt.Time,
			Currency:    row.Currency.String, 
		})
	}
	return list, nil
}

func (r *BookingRepo) UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error {
	bid, _ := uuid.Parse(id)
	return db.New(r.pool).UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
		ID:     pgtype.UUID{Bytes: bid, Valid: true},
		Status: pgtype.Text{String: string(status), Valid: true},
	})
}

func generateResCode() string {
	return fmt.Sprintf("RES-%d", time.Now().UnixNano()%100000)
}

func uuidFromString(s string) [16]byte {
	var b [16]byte
	copy(b[:], s) 
	return b
}

func numericFromFloat(f float64) pgtype.Numeric {
	n := pgtype.Numeric{}
	n.Scan(fmt.Sprintf("%f", f))
	return n
}
