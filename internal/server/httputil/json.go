package httputil

import (
	"github.com/labstack/echo/v4"
)

type OK struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data,omitempty"`
}

type Err struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func RespondOK(c echo.Context, data interface{}) error {
	return c.JSON(200, OK{OK: true, Data: data})
}

func RespondError(c echo.Context, code int, msg string) error {
	return c.JSON(code, Err{OK: false, Error: msg})
}
