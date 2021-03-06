package handler

import (
	"github.com/labstack/echo/v4"
	_ "github.com/mehditeymorian/etefagh/docs"
	"github.com/mehditeymorian/etefagh/internal/model"
	"github.com/mehditeymorian/etefagh/internal/redis"
	"github.com/mehditeymorian/etefagh/internal/request"
	"github.com/mehditeymorian/etefagh/internal/stan"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Event events handler struct
// provide database connection
type Event struct {
	Store  store.Event
	Logger *zap.Logger
	Tracer trace.Tracer
	Stan   stan.Stan
	Redis  redis.Redis
}

// RetrieveAll godoc
// @Summary retrieve all events
// @Description retrieves all events
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {array} model.Event
// @Router /api/v1/events [get]
func (e Event) RetrieveAll(c echo.Context) error {
	ctx, span := e.Tracer.Start(c.Request().Context(), "handler.events.RetrieveAll")
	defer span.End()

	// retrieve all events
	// WARNING!!!!!! require limitation
	all, err := e.Store.RetrieveAll(ctx)
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
		zap.String("method", c.Request().Method),
		zap.String("path", c.Request().RequestURI),
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
	ctx, span := e.Tracer.Start(c.Request().Context(), "handler.events.Retrieve")
	defer span.End()

	eventId := c.Param("event_id")

	retrieve, err := e.Store.Retrieve(ctx, eventId)
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
		zap.String("method", c.Request().Method),
		zap.String("path", c.Request().RequestURI),
	)
	return c.JSON(http.StatusOK, retrieve)
}

// Create godoc
// @Summary create an event
// @Description creates an event
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {object} model.IdResponse
// @Router /api/v1/events [post]
func (e Event) Create(c echo.Context) error {
	ctx, span := e.Tracer.Start(c.Request().Context(), "handler.events.Create")
	defer span.End()

	var publish bool
	var publishType string
	var subject string
	echo.QueryParamsBinder(c).
		Bool("publish", &publish).
		String("publish_type", &publishType).
		String("subject", &subject)

	var input request.Event
	// read body
	if err := c.Bind(&input); err != nil {
		e.Logger.Warn("failed to bind request body",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusBadRequest),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate input
	if err := input.Validate(); err != nil {
		e.Logger.Warn("failed to validate request body",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
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
	id, err := e.Store.Create(ctx, event)
	if err != nil {
		e.Logger.Warn("failed to create event",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	e.Logger.Info("Event Created",
		zap.String("method", c.Request().Method),
		zap.String("path", c.Request().RequestURI),
		zap.Object("input", input),
	)

	// publish event
	if publish {
		err := e.Stan.Publish(ctx, stan.PublishType(publishType), subject, event)

		if err != nil {
			e.Logger.Warn("failed to publish event",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().RequestURI),
				zap.Int("status", http.StatusInternalServerError),
				zap.Error(err),
			)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		e.Logger.Info("event published", zap.Object("event", input))
	}

	// return a struct with id field
	return c.JSON(http.StatusOK, model.IdResponse{Id: id.(primitive.ObjectID).Hex()})
}

// Delete godoc
// @Summary delete an event
// @Description deletes an event by id
// @Tags Event
// @Accept */*
// @Produce json
// @Success 201 {body} struct{}
// @Router /api/v1/events/:event_id [delete]
func (e Event) Delete(c echo.Context) error {
	ctx, span := e.Tracer.Start(c.Request().Context(), "handler.events.Delete")
	defer span.End()

	eventId := c.Param("event_id")

	// delete event
	if err := e.Store.Delete(ctx, eventId); err != nil {
		e.Logger.Warn("failed to delete event with the given id",
			zap.String("event_id", eventId),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	e.Logger.Info("Event Deleted",
		zap.String("event_id", eventId),
		zap.String("method", c.Request().Method),
		zap.String("path", c.Request().RequestURI),
	)
	return c.JSON(http.StatusNoContent, "")
}

// GetPublishStateByEventId godoc
// @Summary get publish state of event
// @Description get publish state of event with asynchronous publish type by event id
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {object} model.StateResponse
// @Router /api/v1/events/:event_id/publish-state [get]
func (e Event) GetPublishStateByEventId(c echo.Context) error {
	ctx, span := e.Tracer.Start(c.Request().Context(), "handler.events.GetPublishStateByEventId")
	defer span.End()

	eventId := c.Param("event_id")

	retrieve, err := e.Store.Retrieve(ctx, eventId)
	if err != nil {
		e.Logger.Warn("failed to retrieve event from db",
			zap.String("event_id", eventId),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// return not found if no event founded
	if retrieve == nil {
		e.Logger.Warn("no event with the given id found",
			zap.String("event_id", eventId),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return c.JSON(http.StatusNotFound, "")
	}

	if retrieve.AckId == "" {
		e.Logger.Info("event published synchronously: no publish state",
			zap.String("event_id", eventId),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return c.JSON(http.StatusOK, model.StateResponse{Description: "event published synchronously"})
	}

	// get publish state
	state, err := e.Redis.GetEventState(ctx, retrieve.AckId)
	if err != nil {
		e.Logger.Warn("failed to get published state",
			zap.String("event_id", eventId),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.StateResponse{State: state})
}

// GetPublishStateByAckId godoc
// @Summary get publish state of event
// @Description get publish state of event with asynchronous publish type by ack id
// @Tags Event
// @Accept */*
// @Produce json
// @Success 200 {object} model.StateResponse
// @Router /api/v1/events/publish-state/:ack_id [get]
func (e Event) GetPublishStateByAckId(c echo.Context) error {
	ctx, span := e.Tracer.Start(c.Request().Context(), "handler.events.GetPublishStateByAckId")
	defer span.End()

	ackId := c.Param("ack_id")

	// get publish state
	state, err := e.Redis.GetEventState(ctx, ackId)
	if err != nil {
		e.Logger.Warn("failed to get published state",
			zap.String("method", c.Request().Method),
			zap.String("path", c.Request().RequestURI),
			zap.Int("status", http.StatusInternalServerError),
			zap.Error(err),
		)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.StateResponse{State: state})
}

// Register register events endpoints on the HTTP server
func (e Event) Register(group *echo.Group) {
	group.GET("/events", e.RetrieveAll)
	group.GET("/events/:event_id", e.Retrieve)
	group.POST("/events", e.Create)
	group.DELETE("/events/:event_id", e.Delete)
	group.GET("/events/:event_id/publish-state", e.GetPublishStateByEventId)
	group.GET("/events/publish-state/:ack_id", e.GetPublishStateByAckId)
}
