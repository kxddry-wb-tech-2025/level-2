package handlers

import (
	"calendar/internal/storage"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func EventsForMonth(st Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, err := strconv.ParseInt(c.QueryParam("uid"), 10, 64)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		date, err := time.Parse("2006-01-02", c.QueryParam("date"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		events, err := st.GetMonth(uid, date)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				return c.String(http.StatusServiceUnavailable, err.Error())
			}
			return c.String(http.StatusBadRequest, err.Error())
		}
		js, err := json.Marshal(events)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, js)
	}
}
