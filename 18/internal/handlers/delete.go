package handlers

import (
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
			return c.String(http.StatusBadRequest, err.Error())
		}

		var s deleteStruct
		if err = json.Unmarshal(data, &s); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err = c.Validate(s); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err = st.Delete(s.ID); err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				return c.String(http.StatusServiceUnavailable, "event not found")
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, fmt.Sprintf("event with id %d deleted", s.ID))
	}
}
