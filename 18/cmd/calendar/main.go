package main

import (
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/storage"
	"fmt"
	"net/http"

	"calendar/internal/validator"

	v10 "github.com/go-playground/validator/v10"
	initCfg "github.com/kxddry/go-utils/pkg/config"
	initLog "github.com/kxddry/go-utils/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var cfg config.Config
	initCfg.MustParseConfig(&cfg)
	log := initLog.SetupLogger(cfg.Env)

	st := storage.New()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST},
		AllowHeaders: []string{},
	}))

	e.Use(middleware.BodyLimit("1M"))
	e.Validator = validator.New(v10.New())

	// POST
	e.POST("/create_event", handlers.CreateEvent(st))
	e.POST("/update_event", handlers.UpdateEvent(st))
	e.POST("/delete_event", handlers.DeleteEvent(st))

	// GET
	e.GET("/events_for_day", handlers.EventsForDay(st))
	e.GET("/events_for_week", handlers.EventsForWeek(st))
	e.GET("/events_for_month", handlers.EventsForMonth(st))

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      e,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error(err.Error())
	}
}
