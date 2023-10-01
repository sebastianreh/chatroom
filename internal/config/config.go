package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		ProjectName    string `default:"chatroom"`
		ProjectVersion string `envconfig:"PROJECT_VERSION" default:"0.0.1"`
		Port           string `envconfig:"PORT" default:"8000" required:"true"`
		Prefix         string `envconfig:"PREFIX" default:"/chatroom"`
		MongoDB        struct {
			Collections struct {
				Users string `envconfig:"USERS_COLLECTION" default:"users"`
				Rooms string `envconfig:"ROOMS_COLLECTION" default:"rooms"`
			}
			Database string `envconfig:"MONGODB_DATABASE" default:"rooms"`
			URI      string `envconfig:"MONGODB_URI" default:"mongodb://localhost:27018"`
		}
	}
)

var (
	Configs Config
)

func NewConfig() Config {
	if err := envconfig.Process("", &Configs); err != nil {
		panic(err.Error())
	}

	return Configs
}
