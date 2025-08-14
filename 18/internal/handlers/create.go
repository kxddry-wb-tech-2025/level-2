package handlers

import (
	"calendar/internal/models"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateEvent(st Storage) echo.HandlerFunc {
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

		out := st.Create(event.UserID, event.Date, event.Text)

		return c.JSON(http.StatusOK, out)
	}
}
