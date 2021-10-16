package handler

import (
	"github.com/labstack/echo/v4"
	_ "github.com/mehditeymorian/etefagh/docs"
	"github.com/mehditeymorian/etefagh/internal/model"
	"github.com/mehditeymorian/etefagh/internal/request"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Event events handler struct
// provide database connection
type Event struct {
	Store  store.Event
	Logger *zap.Logger
}

// RetrieveAll godoc
// @Summary retrieve all events
// @Description retrieves all events
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {object} []model.Event
// @Router /api/v1/events [get]
func (e Event) RetrieveAll(c echo.Context) error {

	// retrieve all events
	// WARNING!!!!!! require limitation
	all, err := e.Store.RetrieveAll(c.Request().Context())
	if err != nil {
		e.Logger.Warn("failed to retrieve all events from db",
			zap.String("path", "/api/v1/events"),
			zap.String("method", "get"),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// set result to an empty array if no event is presents
	if all == nil {
		all = []model.Event{}
	}

	e.Logger.Info("Events Retrieved",
		zap.String("path", "/api/v1/events"),
		zap.String("method", "get"),
		zap.Int("count", len(all)),
	)
	return c.JSON(http.StatusOK, all)
}

// Retrieve godoc
// @Summary retrieve an event
// @Description retrieves an event by id
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {object} model.Event
// @Router /api/v1/events/:event_id [get]
func (e Event) Retrieve(c echo.Context) error {
	eventId := c.Param("event_id")

	retrieve, err := e.Store.Retrieve(c.Request().Context(), eventId)
	if err != nil {
		e.Logger.Warn("failed to retrieve an event from db",
			zap.String("event_id", eventId),
			zap.String("path", "/api/v1/events/:event_id"),
			zap.String("method", "get"),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// return not found if no event founded
	if retrieve == nil {
		e.Logger.Warn("no event with the given id found",
			zap.String("event_id", eventId),
			zap.String("path", "/api/v1/events/:event_id"),
			zap.String("method", "get"),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return c.JSON(http.StatusNotFound, "")
	}

	e.Logger.Info("Event Retrieved",
		zap.String("path", "/api/v1/events/:event_id"),
		zap.String("method", "get"),
	)
	return c.JSON(http.StatusOK, retrieve)
}

// Create godoc
// @Summary create an event
// @Description creates an event
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {string} "HEX ID"
// @Router /api/v1/events [post]
func (e Event) Create(c echo.Context) error {

	var input request.Event
	// read body
	if err := c.Bind(&input); err != nil {
		e.Logger.Warn("failed to bind request body",
			zap.String("path", "/api/v1/events"),
			zap.String("method", "post"),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate input
	if err := input.Validate(); err != nil {
		e.Logger.Warn("failed to validate request body",
			zap.String("path", "/api/v1/events"),
			zap.String("method", "post"),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// recreate event to avoid passing fields that should not be bindable
	event := model.Event{
		EventType:   input.EventType,
		Description: input.Description,
		Priority:    input.Priority,
		Payload:     input.Payload,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		Id:          primitive.NewObjectID(),
	}

	// create event
	id, err := e.Store.Create(c.Request().Context(), event)
	if err != nil {
		e.Logger.Warn("failed to create event",
			zap.String("path", "/api/v1/events"),
			zap.String("method", "post"),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// publish event

	e.Logger.Info("Event Created",
		zap.String("path", "/api/v1/events"),
		zap.String("method", "post"),
		zap.Object("input", input),
	)

	// return a struct with id field
	return c.JSON(http.StatusOK, struct {
		Id string `json:"id"`
	}{Id: id.(primitive.ObjectID).Hex()})
}

// Delete godoc
// @Summary delete an event
// @Description delete an event by id
// @Tags Event
// @Accept */*
// @Produce json
// @Success 201 {string} ""
// @Router /api/v1/events/:event_id [delete]
func (e Event) Delete(c echo.Context) error {
	eventId := c.Param("event_id")

	// delete event
	if err := e.Store.Delete(c.Request().Context(), eventId); err != nil {
		e.Logger.Warn("failed to delete event with the given id",
			zap.String("event_id", eventId),
			zap.String("path", "/api/v1/events"),
			zap.String("method", "post"),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	e.Logger.Info("Event Deleted",
		zap.String("event_id", eventId),
		zap.String("path", "/api/v1/events"),
		zap.String("method", "post"),
	)
	return c.JSON(http.StatusNoContent, "")
}

// Register register events endpoints on the HTTP server
func (e Event) Register(group *echo.Group) {
	group.GET("/events", e.RetrieveAll)
	group.GET("/events/:event_id", e.Retrieve)
	group.POST("/events", e.Create)
	group.DELETE("/events/:event_id", e.Delete)
}
