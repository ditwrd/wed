package server

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware creates a middleware for logging requests
func RequestLoggerMiddleware() echo.MiddlewareFunc {
	logger, _ := zap.NewProduction()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if !strings.Contains(v.URI, "static") {
				logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.String("RemoteIP", v.RemoteIP),
					zap.String("RequestID", v.RequestID),
				)
			}
			return nil
		},
	})
}
