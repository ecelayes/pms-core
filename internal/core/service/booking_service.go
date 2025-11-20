package service

import (
	"context"
	"time"

	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
)

type BookingService struct {
	repo ports.BookingRepository
}

func NewBookingService(repo ports.BookingRepository) ports.BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Create(ctx context.Context, req domain.BookingRequest) (string, error) {
	if len(req.Items) == 0 {
		return "", domain.ErrEmptyBooking
	}
	
	now := time.Now().Truncate(24 * time.Hour)

	for _, item := range req.Items {
		in, err := time.Parse("2006-01-02", item.CheckIn)
		if err != nil {
			return "", domain.ErrInvalidDateFormat
		}

		out, err := time.Parse("2006-01-02", item.CheckOut)
		if err != nil {
			return "", domain.ErrInvalidDateFormat
		}

		if in.Before(now) {
			return "", domain.ErrPastDate
		}

		if !out.After(in) {
			return "", domain.ErrInvalidDateRange
		}
	}
	
	return s.repo.CreateBookingAtomic(ctx, req)
}

func (s *BookingService) Get(ctx context.Context, id string) (*domain.Booking, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookingService) List(ctx context.Context, hotelID string) ([]domain.Booking, error) {
	return s.repo.List(ctx, hotelID)
}

func (s *BookingService) Cancel(ctx context.Context, id string) error {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if b.Status == domain.BookingStatusCancelled {
		return nil
	}

	return s.repo.CancelBookingAtomic(ctx, id)
}
