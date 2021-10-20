package telemetry

type Config struct {
	Trace `yaml:"trace"`
}

type Trace struct {
	Enabled bool   `yaml:"enabled"`
	Url     string `yaml:"url"`
}
