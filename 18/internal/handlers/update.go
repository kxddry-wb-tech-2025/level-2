package handlers

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/kxddry/go-utils/pkg/logger/handlers/sl"
	"github.com/labstack/echo/v4"
)

func UpdateEvent(st Storage, log *slog.Logger) echo.HandlerFunc {
	const op = "handlers.UpdateEvent"
	log = log.With(slog.String("op", op))
	return func(c echo.Context) error {
		body := c.Request().Body
		defer func() { _ = body.Close() }()
		data, err := io.ReadAll(body)
		if err != nil {
			log.Debug("case 1", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		var event models.Event
		if err := json.Unmarshal(data, &event); err != nil {
			log.Debug("case 2", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err := c.Validate(event); err != nil {
			log.Debug("case 3", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if event.ID == 0 {
			log.Debug("no id / zero id", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: errors.New("invalid event id")})
		}

		out, err := st.Update(event.ID, event.Date, event.Text)
		if err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				log.Debug("case 4", sl.Err(err))
				return c.JSON(http.StatusServiceUnavailable, models.Response{Error: errors.New("event not found")})
			}
			log.Debug("case 5", sl.Err(err))
			return c.JSON(http.StatusInternalServerError, models.Response{Error: err})
		}

		return c.JSON(http.StatusOK, models.Response{Result: out})
	}
}
