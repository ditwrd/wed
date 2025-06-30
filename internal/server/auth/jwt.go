package auth

import (
	"time"

	"github.com/ditwrd/wed/internal/server/httputil"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// ===============================================
// JWT Utilities
// ===============================================

// JWTSecret is the secret key for signing JWT tokens
var JWTSecret []byte

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AdminAuthMiddleware checks for valid JWT token
func AdminAuthMiddleware() echo.MiddlewareFunc {
	JWTSecret = []byte(viper.GetString("app.jwt_secret"))

	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(Claims) },
		SigningKey:    JWTSecret,
		TokenLookup:   "header:Authorization:Bearer ,cookie:token",
		ErrorHandler: func(c echo.Context, err error) error {
			if c.Path() == "/api/v1/rsvps" || c.Path() == "/api/v1/stats" {
				return httputil.RespondError(c, 401, "unauthorized")
			}
			return c.Redirect(303, "/admin/login")
		},
		SigningMethod: "HS256",
	})
}

// CreateAdminToken creates a new JWT token for admin user
func CreateAdminToken() (string, error) {
	JWTSecret = []byte(viper.GetString("app.jwt_secret"))

	// Create JWT token
	now := time.Now()
	expirationTime := now.Add(2 * time.Hour)
	claims := &Claims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.Add(-1 * time.Minute)),
			Issuer:    "wed-admin",
			Audience:  jwt.ClaimStrings{"wed-admin"},
			Subject:   "admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}
