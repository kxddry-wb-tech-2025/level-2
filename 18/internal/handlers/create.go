package handlers

import (
	"calendar/internal/models"
	"encoding/json"
	"errors"
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
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		var event models.Event
		if err := json.Unmarshal(data, &event); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err := c.Validate(event); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if event.ID != 0 {
			return c.JSON(http.StatusBadRequest,
				models.Response{
					Error: errors.New("you cannot set your own id, its assigned by the server"),
				})
		}

		out := st.Create(event.UserID, event.Date, event.Text)

		return c.JSON(http.StatusOK, models.Response{Result: out})
	}
}
