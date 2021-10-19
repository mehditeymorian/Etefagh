# Etefagh
Publish events to Nats Streaming(STAN).
![Structure](https://github.com/mehditeymorian/Etefagh/blob/master/assets/ETEFAGH.jpg)

### Event Model
```golang
// Event model
type Event struct {
	EventType   string `json:"event_type"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Payload     string `json:"payload"`
}
```

