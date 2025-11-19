package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	
	httpHandler "github.com/ecelayes/pms-core/internal/adapter/handler/http"
)

type Server struct {
	Echo *echo.Echo
}

func New() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	return &Server{Echo: e}
}

func (s *Server) RegisterRoutes(avHandler *httpHandler.AvailabilityHandler) {
	v1 := s.Echo.Group("/api/v1")

	v1.GET("/health", func(c echo.Context) error { return c.JSON(200, "ok") })
	
	v1.GET("/availability", avHandler.Get)
}

func (s *Server) Start(port string) error {
	return s.Echo.Start(":" + port)
}
