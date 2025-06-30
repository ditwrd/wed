package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/ditwrd/wed/internal/web/page"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// ===============================================
// Page Handlers
// ===============================================

// AdminLogin shows the login form
func AdminLogin(c echo.Context) error {
	return Render(c, http.StatusOK, page.AdminLogin(page.AdminLoginProps{}))
}

// ===============================================
// Form Handlers
// ===============================================

// AdminLoginPost handles the login form submission
func AdminLoginPost(c echo.Context) error {
	// Get the password from the form
	password := c.FormValue("password")
	if password == "" {
		return Render(c, http.StatusOK, page.AdminLogin(page.AdminLoginProps{
			Error: "Password is required",
		}))
	}

	// Get the stored password hash from config
	storedHash := viper.GetString("app.admin_password_hash")
	if storedHash == "" {
		return Render(c, http.StatusOK, page.AdminLogin(page.AdminLoginProps{
			Error: "Admin password not configured",
		}))
	}

	// Compare the password with the stored hash
	fmt.Println(password)
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return Render(c, http.StatusOK, page.AdminLogin(page.AdminLoginProps{
			Error: "Invalid password",
		}))
	}

	// Create JWT token
	tokenString, err := CreateAdminToken()
	if err != nil {
		return Render(c, http.StatusOK, page.AdminLogin(page.AdminLoginProps{
			Error: "Failed to generate token",
		}))
	}

	// Set the token as a cookie
	expirationTime := time.Now().Add(2 * time.Hour)
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false, // Set to true for production with HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	c.SetCookie(cookie)

	// Redirect to admin dashboard
	return c.Redirect(http.StatusSeeOther, "/admin")
}

// Render is a helper function to render templ components
func Render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}
