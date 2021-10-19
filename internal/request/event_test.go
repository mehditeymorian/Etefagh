package request_test

import (
	"github.com/mehditeymorian/etefagh/internal/request"
	"testing"
)

func TestEventValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		eventType   string
		description string
		priority    int
		payload     string
		isValid     bool
	}{
		{
			eventType:   "",
			description: "desc1",
			priority:    5,
			payload:     "{}",
			isValid:     false,
		},
		{
			eventType:   "test-type",
			description: "desc1",
			priority:    0,
			payload:     "{}",
			isValid:     false,
		},
		{
			eventType:   "test-type",
			description: "desc1",
			priority:    12,
			payload:     "{}",
			isValid:     false,
		},
		{
			eventType:   "test-type",
			description: "desc1",
			priority:    1,
			payload:     "{hello}",
			isValid:     false,
		},
		{
			eventType:   "event1",
			description: "this is a description",
			priority:    4,
			payload:     `{"name":"mehdi"}`,
			isValid:     true,
		},
	}

	for _, c := range cases {

		rq := request.Event{
			EventType:   c.eventType,
			Description: c.description,
			Priority:    c.priority,
			Payload:     c.payload,
		}

		err := rq.Validate()

		// request is valid but validation has error
		if c.isValid && err != nil {
			t.Fatalf("valid request %+v has error %s", rq, err)
		}

		// request is invalid but validation has no error
		if !c.isValid && err == nil {
			t.Fatalf("invalid request %+v has no error %s", rq, err)
		}
	}

}
