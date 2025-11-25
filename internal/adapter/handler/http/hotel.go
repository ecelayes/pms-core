package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/ports"
	"github.com/ecelayes/pms-core/pkg/response"
)

type HotelHandler struct {
	repo ports.HotelRepository
}

func NewHotelHandler(repo ports.HotelRepository) *HotelHandler {
	return &HotelHandler{repo: repo}
}

func (h *HotelHandler) Create(c echo.Context) error {
	var req domain.CreateHotelRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err)
	}

	userID := c.Get("user_id").(string)

	id, err := h.repo.CreateHotel(c.Request().Context(), userID, req)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}

	return response.Success(c, http.StatusCreated, map[string]string{"hotel_id": id})
}

func (h *HotelHandler) ListMine(c echo.Context) error {
	userID := c.Get("user_id").(string)
	
	hotels, err := h.repo.ListByOwner(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}
	return response.Success(c, http.StatusOK, hotels)
}
