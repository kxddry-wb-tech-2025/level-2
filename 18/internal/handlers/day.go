package handlers

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func EventsForDay(st Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, err := strconv.ParseInt(c.QueryParam("uid"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		date, err := time.Parse("2006-01-02", c.QueryParam("date"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		events, err := st.GetDay(uid, date)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				return c.JSON(http.StatusServiceUnavailable, models.Response{Error: err})
			}
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}
		return c.JSON(http.StatusOK, models.Response{Result: events})
	}
}
