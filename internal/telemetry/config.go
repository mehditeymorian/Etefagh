package telemetry

type Config struct {
	Trace
}

type Trace struct {
	Enabled bool
	Url     string
}
