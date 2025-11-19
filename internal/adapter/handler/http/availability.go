package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ecelayes/pms-core/internal/core/ports"
	"github.com/ecelayes/pms-core/pkg/response"
)

type AvailabilityHandler struct {
	repo ports.AvailabilityRepository
}

func NewAvailabilityHandler(repo ports.AvailabilityRepository) *AvailabilityHandler {
	return &AvailabilityHandler{repo: repo}
}

func (h *AvailabilityHandler) Get(c echo.Context) error {
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	if start == "" || end == "" {
		return response.Error(c, http.StatusBadRequest, echo.ErrBadRequest)
	}

	data, err := h.repo.GetAvailability(c.Request().Context(), start, end)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, err)
	}

	return response.Success(c, http.StatusOK, data)
}
