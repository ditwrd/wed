package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/ditwrd/wed/internal/server/page"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/ratelimit"
)

func NewServer() *echo.Echo {
	e := echo.New()
	// Pre() middleware order
	if !viper.GetBool("app.dev") {
		e.Pre(middleware.HTTPSRedirect())
	}
	e.Pre(middleware.NonWWWRedirect())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.MethodOverride())
	// Use() middleware stack
	// e.Use(middleware.Recover())
	e.Use(RequestLoggerMiddleware())
	e.Use(middleware.Secure())
	// Configurable CORS via viper
	corsOrigins := viper.GetStringSlice("app.cors.origins")
	corsMethods := viper.GetStringSlice("app.cors.methods")
	corsHeaders := viper.GetStringSlice("app.cors.headers")
	corsCredentials := viper.GetBool("app.cors.credentials")
	corsMaxAge := viper.GetInt("app.cors.max_age")
	corsCfg := middleware.CORSConfig{
		AllowOrigins:     corsOrigins,
		AllowMethods:     corsMethods,
		AllowHeaders:     corsHeaders,
		AllowCredentials: corsCredentials,
	}
	if corsMaxAge > 0 {
		corsCfg.MaxAge = corsMaxAge
	}
	if len(corsCfg.AllowOrigins) == 0 {
		corsCfg.AllowOrigins = []string{"*"}
	}
	if len(corsCfg.AllowMethods) == 0 {
		corsCfg.AllowMethods = []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		}
	}
	if len(corsCfg.AllowHeaders) == 0 {
		corsCfg.AllowHeaders = []string{"*"}
	}
	e.Use(middleware.CORSWithConfig(corsCfg))

	// Global rate limiter
	if !viper.GetBool("app.dev") {
		rps := viper.GetInt("app.rate_limit_rps")
		if rps <= 0 {
			rps = 5
		}
		rl := ratelimit.New(rps)
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				_ = rl.Take()
				return next(c)
			}
		})
	}

	e.GET("/static/*", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=31536000, immutable, stale-while-revalidate=86400")
		return nil
	}).Name = "static"

	e.Use(
		middleware.GzipWithConfig(
			middleware.GzipConfig{
				Skipper: func(c echo.Context) bool { return strings.Contains(c.Path(), "metrics") },
			},
		),
	)
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: 30 * 1e9}))
	// Prometheus metrics
	// e.Use(echoprometheus.NewMiddleware("wed"))
	// e.GET("/metrics", echoprometheus.NewHandler())
	// Health endpoint
	e.GET("/healthz", func(c echo.Context) error { return c.NoContent(200) })
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := 500
		msg := "Internal Server Error"
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if m, ok := he.Message.(string); ok {
				msg = m
			} else if he.Message != nil {
				if em, ok := he.Message.(error); ok {
					msg = em.Error()
				}
			}
		}
		_ = page.RenderError(c, code, msg)
	}
	return e
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
			// Request logger configured in RequestLoggerMiddleware()

			go startServer(e)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}
