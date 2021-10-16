package server

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehditeymorian/etefagh/internal/config"
	"github.com/mehditeymorian/etefagh/internal/db"
	handler "github.com/mehditeymorian/etefagh/internal/handler/event"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Main(config config.Config) {

	// create HTTP server
	app := echo.New()

	// create a database connection
	database, err := db.Connect(config.DB)
	if err != nil {
		panic(fmt.Errorf("database initialization failed %w", err))
	}

	// register events endpoints
	handler.Event{Store: store.NewMongoEvent(database)}.Register(app.Group("/api/v1"))

	// start HTTP Server
	if err := app.Start(fmt.Sprintf(":%s", config.Api.Port)); !errors.Is(err, http.ErrServerClosed) {
		panic("echo initiation failed")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
