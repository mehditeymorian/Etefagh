package redis

type Config struct {
	Address  string // host:port
	Password string // "" means no password
	DB       int    // 0 means default DB
}
