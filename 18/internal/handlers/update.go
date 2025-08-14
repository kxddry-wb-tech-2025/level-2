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
			return c.String(http.StatusBadRequest, err.Error())
		}

		var event models.Event
		if err := json.Unmarshal(data, &event); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(event); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if event.ID == 0 {
			return c.String(http.StatusBadRequest, "invalid event id")
		}

		out, err := st.Update(event.ID, event.Date, event.Text)
		if err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				return c.String(http.StatusServiceUnavailable, "event not found")
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, out)
	}
}
