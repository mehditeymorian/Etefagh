package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/mehditeymorian/etefagh/internal/model"
	"github.com/mehditeymorian/etefagh/internal/request"
	store "github.com/mehditeymorian/etefagh/internal/store/event"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
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

	var input request.Event
	// read body
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// validate input
	if err := input.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// recreate event to avoid passing fields that should not be bindable
	event := model.Event{
		EventType:   input.EventType,
		Description: input.Description,
		Priority:    input.Priority,
		Payload:     input.Payload,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		Id:          primitive.ObjectID{},
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
