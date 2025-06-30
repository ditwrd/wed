package server

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewServer() *echo.Echo {
	return echo.New()
}

func startServer(e *echo.Echo) {
	bind := viper.GetString("app.bind")
	port := viper.GetString("app.port")
	address := strings.Join([]string{bind, port}, ":")
	err := e.Start(address)
	if err != nil {
		panic(err)
	}
}

func RegisterHooks(lifecycle fx.Lifecycle, e *echo.Echo) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger, _ := zap.NewProduction()
			e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
				LogURI:    true,
				LogStatus: true,
				LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
					if !strings.Contains(v.URI, "static") {
						logger.Info("request",
							zap.String("URI", v.URI),
							zap.Int("status", v.Status),
						)
					}
					return nil
				},
			}))

			go startServer(e)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}