package handlers

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/kxddry/go-utils/pkg/logger/handlers/sl"
	"github.com/labstack/echo/v4"
)

type deleteStruct struct {
	ID int64 `json:"id" validate:"required"`
}

func DeleteEvent(st Storage, log *slog.Logger) echo.HandlerFunc {
	const op = "handlers.DeleteEvent"
	log = log.With(slog.String("op", op))
	return func(c echo.Context) error {
		body := c.Request().Body
		defer func() { _ = body.Close() }()

		data, err := io.ReadAll(body)
		if err != nil {
			log.Debug("case 1", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		var s deleteStruct
		if err = json.Unmarshal(data, &s); err != nil {
			log.Debug("case 2", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err = c.Validate(s); err != nil {
			log.Debug("case 3", sl.Err(err))
			return c.JSON(http.StatusBadRequest, models.Response{Error: err})
		}

		if err = st.Delete(s.ID); err != nil {
			if errors.Is(err, storage.ErrEventNotFound) {
				log.Debug("case 4", sl.Err(err))
				return c.JSON(http.StatusServiceUnavailable, "event not found")
			}
			log.Debug("case 5", sl.Err(err))
			return c.JSON(http.StatusInternalServerError, models.Response{Error: err})
		}

		return c.JSON(http.StatusOK, models.Response{Result: fmt.Sprintf("event with id %d deleted", s.ID)})
	}
}
