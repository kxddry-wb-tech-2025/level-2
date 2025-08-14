package handlers

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateEvent(st Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := c.Request().Body
		defer func() { _ = body.Close() }()
		data, err := io.ReadAll(body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		var event models.Event
		if err := json.Unmarshal(data, &event); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err := c.Validate(event); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if event.ID == 0 {
			return c.JSON(http.StatusBadRequest, models.Response{Error: errors.New("invalid event id")})
		}

		out, err := st.Update(event.ID, event.Date, event.Text)
		if err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				return c.JSON(http.StatusServiceUnavailable, models.Response{Error: errors.New("event not found")})
			}
			return c.JSON(http.StatusInternalServerError, models.Response{Error: err})
		}

		return c.JSON(http.StatusOK, models.Response{Result: out})
	}
}
