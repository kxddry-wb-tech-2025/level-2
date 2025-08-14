package handlers

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type deleteStruct struct {
	ID int64 `json:"id" validate:"required"`
}

func DeleteEvent(st Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		body := c.Request().Body
		defer func() { _ = body.Close() }()

		data, err := io.ReadAll(body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		var s deleteStruct
		if err = json.Unmarshal(data, &s); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err = c.Validate(s); err != nil {
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err = st.Delete(s.ID); err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				return c.JSON(http.StatusServiceUnavailable, "event not found")
			}
			return c.JSON(http.StatusInternalServerError, models.Response{Error: err})
		}

		return c.JSON(http.StatusOK, models.Response{Result: fmt.Sprintf("event with id %d deleted", s.ID)})
	}
}
