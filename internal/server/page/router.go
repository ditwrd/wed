package page

import (
	"net/http"

	"github.com/ditwrd/wed/internal/server/auth"
	"github.com/labstack/echo/v4"
)

func Router(e *echo.Echo, h *Handler) {
	e.GET("/static/*", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable, stale-while-revalidate=86400")
		http.StripPrefix("/static/", http.FileServer(h.Static())).ServeHTTP(c.Response(), c.Request())
		return nil
	}).Name = "static"
	e.GET("/", h.Home).Name = "home"
	e.GET("/p/:name", h.PersonalizedHome).Name = "home_personal"
	e.GET("/g/:group", h.GroupHome).Name = "home_group"
	e.POST("/rsvp", h.HandleRSVP).Name = "rsvp_create"

	// Admin login routes (no authentication required)
	e.GET("/admin/login", auth.AdminLogin).Name = "admin_login"
	e.POST("/admin/auth/login", auth.AdminLoginPost).Name = "admin_login_post"

	// Admin routes with authentication
	admin := e.Group("/admin", auth.AdminAuthMiddleware())
	admin.GET("", h.AdminDashboard).Name = "admin_dashboard"
	admin.GET("/logout", auth.AdminLogout).Name = "admin_logout"

	// Public API v1
	api := e.Group("/api/v1")
	api.GET("/rsvps", h.AdminRSVPList).Name = "api_v1_rsvps"
	api.GET("/stats", h.AdminRSVPStats).Name = "api_v1_stats"
}
