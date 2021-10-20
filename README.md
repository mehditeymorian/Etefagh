# Etefagh
Publishes events to Nats Streaming(STAN) synchornously and asynchronously. Cache events's publish-state using Redis and Store events using MongoDB.
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

### Protobuf
```protobuf
message Event {
  string eventType = 1;
  string description = 2;
  int32 priority = 3;
  string payload = 4;
}
```

### How to Run
Use [docker-compose.yaml](https://github.com/mehditeymorian/Etefagh/blob/master/deployments/docker-compose.yaml) file to run containers.

- Up: `docker-compose -f PATH/deployments/docker-compose.yaml up --build`
- Down: `docker-compose -f PATH/deployments/docker-compose.yaml down`

