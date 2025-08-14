package handlers

import (
	"calendar/internal/models"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/kxddry/go-utils/pkg/logger/handlers/sl"
	"github.com/labstack/echo/v4"
)

func CreateEvent(st Storage, log *slog.Logger) echo.HandlerFunc {
	const op = "handlers.CreateEvent"
	log = log.With(slog.String("op", op))
	return func(c echo.Context) error {
		body := c.Request().Body
		defer func() { _ = body.Close() }()

		data, err := io.ReadAll(body)
		if err != nil {
			log.Debug("case 1", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		}

		var event models.Event
		if err := json.Unmarshal(data, &event); err != nil {
			log.Debug("case 2", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		}

		if err := c.Validate(event); err != nil {
			log.Debug("case 3", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		}

		if event.ID != 0 {
			log.Debug("case 4", sl.Err(err))
			return c.JSON(http.StatusBadRequest,
				models.Response{
					Error: "you cannot set your own id, its assigned by the server",
				})
		}

		out := st.Create(event.UserID, event.Date, event.Text)

		return c.JSON(http.StatusOK, models.Response{Result: out})
	}
}
