package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ecelayes/pms-core/internal/core/domain"
	"github.com/ecelayes/pms-core/internal/core/service"
	"github.com/ecelayes/pms-core/pkg/response"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err)
	}

	id, err := h.svc.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return response.Error(c, http.StatusConflict, err)
		}

		return response.Error(c, http.StatusInternalServerError, err)
	}

	return response.Success(c, http.StatusCreated, map[string]string{"user_id": id})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err)
	}

	res, err := h.svc.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return response.Error(c, http.StatusUnauthorized, err)
		}
		
		return response.Error(c, http.StatusInternalServerError, err)
	}

	return response.Success(c, http.StatusOK, res)
}
