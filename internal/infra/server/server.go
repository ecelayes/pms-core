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
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	return &Server{Echo: e}
}

func (s *Server) RegisterRoutes(
	avHandler *httpHandler.AvailabilityHandler,
	bookHandler *httpHandler.BookingHandler,
	authHandler *httpHandler.AuthHandler,
	hotelHandler *httpHandler.HotelHandler,
) {
	v1 := s.Echo.Group("/api/v1")

	v1.GET("/health", func(c echo.Context) error { return c.JSON(200, "ok") })

	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	admin := v1.Group("/admin")
	admin.Use(AuthMiddleware)

	admin.POST("/hotels", hotelHandler.Create)
	admin.GET("/hotels", hotelHandler.ListMine)

	v1.GET("/availability", avHandler.Get)

	v1.POST("/bookings", bookHandler.Create)
  v1.GET("/bookings", bookHandler.List)
  v1.GET("/bookings/:id", bookHandler.Get)
  v1.POST("/bookings/:id/cancel", bookHandler.Cancel)
}

func (s *Server) Start(port string) error {
	return s.Echo.Start(":" + port)
}
