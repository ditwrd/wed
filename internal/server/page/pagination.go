package page

import (
	"fmt"
	"net/http"

	"github.com/ditwrd/wed/internal/server/httputil"
	"github.com/labstack/echo/v4"
)

func ParsePagination(c echo.Context) (int, int, error) {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")
	if limit == "" {
		limit = "10"
	}
	if offset == "" {
		offset = "0"
	}
	limitInt := 10
	if _, err := fmt.Sscanf(limit, "%d", &limitInt); err != nil {
		return 0, 0, err
	}
	offsetInt := 0
	if _, err := fmt.Sscanf(offset, "%d", &offsetInt); err != nil {
		return 0, 0, err
	}
	if limitInt < 1 || limitInt > 100 {
		limitInt = 10
	}
	if offsetInt < 0 {
		offsetInt = 0
	}
	return limitInt, offsetInt, nil
}

func RespondInvalidPagination(c echo.Context) error {
	return httputil.RespondError(c, http.StatusBadRequest, "invalid pagination")
}
