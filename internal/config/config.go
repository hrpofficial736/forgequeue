package config

type Config struct {
	RedisURL    string
	PostgresURL string
}

func LoadConfig() *Config {
	cfg := &Config{
		RedisURL: "localhost:6379",
		PostgresURL: "localhost:5432",
	}
	
	return cfg
}