package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mehditeymorian/etefagh/internal/model"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Event events handler struct
// provide database connection
type Event struct {
	Store store.Event
}

// RetrieveAll events
func (e Event) RetrieveAll(c echo.Context) error {

	// retrieve all events
	// WARNING!!!!!! require limitation
	all, err := e.Store.RetrieveAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// set result to an empty array if no event is presents
	if all == nil {
		all = []model.Event{}
	}

	return c.JSON(http.StatusOK, all)
}

// Retrieve an event by id
func (e Event) Retrieve(c echo.Context) error {
	eventId := c.Param("event_id")

	retrieve, err := e.Store.Retrieve(c.Request().Context(), eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// return not found if no event founded
	if retrieve == nil {
		return c.JSON(http.StatusNotFound, "")
	}

	return c.JSON(http.StatusOK, retrieve)
}

// Create an event
func (e Event) Create(c echo.Context) error {

	var event model.Event
	// read body
	if err := c.Bind(&event); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// recreate event to avoid passing fields that should not be bindable
	event = model.Event{
		EventType:   event.EventType,
		Description: event.Description,
		Priority:    event.Priority,
		Payload:     event.Payload,
		CreatedAt:   0,
		Id:          primitive.ObjectID{},
	}

	// validate input
	if err := event.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// create event
	id, err := e.Store.Create(c.Request().Context(), event)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// publish event

	return c.JSON(http.StatusOK, id)
}

// Delete an event
func (e Event) Delete(c echo.Context) error {
	eventId := c.Param("event_id")

	// delete event
	if err := e.Store.Delete(c.Request().Context(), eventId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusNoContent, "")
}

// Register register events endpoints on the HTTP server
func (e Event) Register(group *echo.Group) {
	group.GET("/events", e.RetrieveAll)
	group.GET("/events/:event_id", e.Retrieve)
	group.POST("/events", e.Create)
	group.DELETE("/events/:event_id", e.Delete)
}
