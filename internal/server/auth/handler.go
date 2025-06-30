package auth

import (
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

// AdminLogout clears the auth cookie and redirects to login
func AdminLogout(c echo.Context) error {
	domain := viper.GetString("app.domain_name")
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Domain:   domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   !viper.GetBool("app.dev"),
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/admin/login")
}

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
	expirationTime := time.Now().Add(1 * time.Hour)
	domain := viper.GetString("app.domain_name")
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   !viper.GetBool("app.dev"),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Domain:   domain,
	}
	c.SetCookie(cookie)

	// Redirect to admin dashboard
	return c.Redirect(http.StatusSeeOther, "/admin")
}

// Render is a local helper to avoid import cycles
func Render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)
	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}
	return c.HTML(statusCode, buf.String())
}
