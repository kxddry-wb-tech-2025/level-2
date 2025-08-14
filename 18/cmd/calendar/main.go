package main

import (
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/storage"

	"calendar/internal/validator"

	v10 "github.com/go-playground/validator/v10"
	initCfg "github.com/kxddry/go-utils/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	var cfg config.Config
	initCfg.MustParseConfig(&cfg)

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

	e.POST("/create_event", handlers.CreateEvent(st))
	e.POST("/update_event", handlers.UpdateEvent(st))
}
