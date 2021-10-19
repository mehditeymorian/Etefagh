package eventpb

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/mehditeymorian/etefagh/internal/model"
)

type MarshalType string

const (
	Json     MarshalType = "Json"
	Protobuf             = "Protobuf"
)

func Marshal(event model.Event, marshalType MarshalType) ([]byte, error) {
	switch marshalType {
	case Protobuf:
		return protoMarshal(event)
	default:
	case Json:
		return json.Marshal(event)
	}

	return nil, nil
}

func protoMarshal(event model.Event) ([]byte, error) {
	protoEvent := Event{
		EventType:   event.EventType,
		Description: event.Description,
		Priority:    int32(event.Priority),
		Payload:     event.Payload,
	}
	return proto.Marshal(&protoEvent)
}
