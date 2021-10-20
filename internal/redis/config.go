package redis

type Config struct {
	Address  string `yaml:"address"`  // host:port
	Password string `yaml:"password"` // "" means no password
	DB       int    `yaml:"db"`       // 0 means default DB
}
