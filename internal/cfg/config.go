package cfg

import "os"

var TestinatorConfig Config

func init() {
	TestinatorConfig = Config{
		AppURL: os.Getenv("APP_URL"),
	}
}

type Config struct {
	AppURL string
}
