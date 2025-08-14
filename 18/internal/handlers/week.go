package handlers

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/kxddry/go-utils/pkg/logger/handlers/sl"
	"github.com/labstack/echo/v4"
)

func EventsForWeek(st Storage, log *slog.Logger) echo.HandlerFunc {
	const op = "handlers.EventsForWeek"
	log = log.With(slog.String("op", op))
	return func(c echo.Context) error {
		uidStr := c.QueryParam("user_id")
		if uidStr == "" {
			log.Debug("case 0")
			return c.JSON(http.StatusBadRequest, models.Response{Error: "user_id is empty"})
		}

		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			log.Debug("case 1", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		}

		date, err := time.Parse("2006-01-02", c.QueryParam("date"))
		if err != nil {
			log.Debug("case 2", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		}

		events, err := st.GetWeek(uid, models.Date(date))
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Debug("case 3", sl.Err(err))
				return c.JSON(http.StatusServiceUnavailable, models.Response{Error: err.Error()})
			}
			log.Debug("case 4", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, models.Response{Result: events})
	}
}
