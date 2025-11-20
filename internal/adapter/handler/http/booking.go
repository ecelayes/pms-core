package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
	"github.com/ecelayes/pms-core/pkg/response"
)

type BookingHandler struct {
	svc ports.BookingService
}

func NewBookingHandler(svc ports.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) Create(c echo.Context) error {
	var req domain.BookingRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err)
	}

	code, err := h.svc.Create(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNoAvailability):
			return response.Error(c, http.StatusConflict, err)
			
		case errors.Is(err, domain.ErrInvalidHotelID),
			 errors.Is(err, domain.ErrEmptyBooking),
			 errors.Is(err, domain.ErrInvalidDateFormat),
			 errors.Is(err, domain.ErrInvalidDateRange),
			 errors.Is(err, domain.ErrPastDate):
			return response.Error(c, http.StatusBadRequest, err)
		}

		return response.Error(c, http.StatusInternalServerError, err)
	}

	return response.Success(c, http.StatusCreated, map[string]string{
		"booking_code": code,
		"status":       "confirmed",
	})
}

func (h *BookingHandler) Get(c echo.Context) error {
	id := c.Param("id")
	booking, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusNotFound, err)
	}
	return response.Success(c, http.StatusOK, booking)
}

func (h *BookingHandler) List(c echo.Context) error {
	hotelID := c.QueryParam("hotel_id")
	if hotelID == "" {
		return response.Error(c, http.StatusBadRequest, errors.New("hotel_id is required"))
	}

	list, err := h.svc.List(c.Request().Context(), hotelID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}
	return response.Success(c, http.StatusOK, list)
}

func (h *BookingHandler) Cancel(c echo.Context) error {
	id := c.Param("id")
	if err := h.svc.Cancel(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}
	return response.Success(c, http.StatusOK, map[string]string{"status": "cancelled"})
}
