package server

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mehditeymorian/etefagh/docs"
	"github.com/mehditeymorian/etefagh/internal/config"
	"github.com/mehditeymorian/etefagh/internal/db"
	handler "github.com/mehditeymorian/etefagh/internal/handler/event"
	"github.com/mehditeymorian/etefagh/internal/redis"
	stan "github.com/mehditeymorian/etefagh/internal/stan"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title Swagger API
// @version 1.0
// @description This is an Event Publisher server.

// @contact.name Mehdi Teymorian
// @contact.url https://www.timurid.ir
// @contact.email mehditeymorian322@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// Main
// @host http://localhost
// @BasePath /api/v1
func Main(config config.Config, logger *zap.Logger, tracer trace.Tracer) {

	// create HTTP server
	app := echo.New()

	// create a database connection
	database, err := db.Connect(config.DB)
	if err != nil {
		err := fmt.Errorf("database initialization failed %w", err)
		logger.Fatal(err.Error())
	}
	mongoStore := store.NewMongoEvent(database, tracer)

	// create redis connection
	redis := redis.Connect(config.Redis, tracer)

	// connect to stan
	stanConn, err := stan.Connect(config.Nats)
	stan := stan.Stan{
		Connection: stanConn,
		Redis:      redis,
		Store:      mongoStore,
		Tracer:     tracer,
	}

	if err != nil {
		err := fmt.Errorf("stan initalization failed: %w", err)
		logger.Fatal(err.Error())
	}

	// register swagger
	app.GET("/swagger/*", echoSwagger.WrapHandler)

	// register events endpoints
	handler.Event{
		Store:  mongoStore,
		Logger: logger,
		Tracer: tracer,
		Stan:   stan,
		Redis:  redis,
	}.Register(app.Group("/api/v1"))

	// start HTTP Server
	if err := app.Start(fmt.Sprintf(":%s", config.Api.Port)); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Echo failed to start", zap.String("port", config.Api.Port))
		panic("Echo failed to start")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
