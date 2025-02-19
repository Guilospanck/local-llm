package pkg

import (
	"base/pkg/domain"
	"base/pkg/handlers"
	"base/pkg/utils"
	"errors"
	"flag"
	"fmt"

	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *domain.Database
var modelName *string

func Server() {
	e := echo.New()

	db = domain.NewDb()
	db.Connect()

	defer db.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	modelName = flag.String("model", "deepseek-r1", "Deep Seek R1")
	flag.Parse()

	// Routes
	e.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "I'm healthy") })
	e.POST("/", func(c echo.Context) error {
		return handlers.CallStream(c, modelName)
	})
	e.POST("/call", func(c echo.Context) error {
		return handlers.CallSimple(c, modelName)
	})
	e.POST("/extract", func(c echo.Context) error {
		return handlers.Extract(c, db, modelName)
	})

	// Start server
	if err := e.Start(fmt.Sprintf(":%d", utils.SERVER_PORT)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
