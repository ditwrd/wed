package page

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(c.Request().Context(), buf); err != nil {
		return err
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.HTML(statusCode, buf.String())
	}
	return c.HTML(statusCode, buf.String())
}

func RenderError(c echo.Context, statusCode int, msg string) error {
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.String(statusCode, msg)
	}
	return c.String(statusCode, msg)
}
