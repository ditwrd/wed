package page

import (
	"net/http"

	"github.com/ditwrd/wed/internal/server/auth"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo, h *Handler) {
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(h.Static()))))
	e.GET("/", h.Home)
	e.GET("/p/:name", h.PersonalizedHome)
	e.GET("/g/:group", h.GroupHome)
	e.POST("/rsvp", h.HandleRSVP)

	// Admin login routes (no authentication required)
	e.GET("/admin/login", auth.AdminLogin)
	e.POST("/admin/auth/login", auth.AdminLoginPost)

	// Admin routes with authentication
	admin := e.Group("/admin", auth.AdminAuthMiddleware())
	admin.GET("", h.AdminDashboard)
	admin.GET("/api/rsvps", h.AdminRSVPList)
	admin.GET("/api/stats", h.AdminRSVPStats)
}
