package cfg

import "os"

var TestinatorConfig Config

func init() {
	TestinatorConfig = Config{
		AppURL: os.Getenv("APP_URL"),
		Postgres: PostgresConfig{
			ConnectionURL: os.Getenv("POSTGRES_CONNECTION_URL"),
		},
	}
}

type Config struct {
	AppURL   string
	Postgres PostgresConfig
}

type PostgresConfig struct {
	ConnectionURL string
}
