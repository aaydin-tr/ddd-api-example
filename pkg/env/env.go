package env

import (
	"log"
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type ENV struct {
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     int    `env:"POSTGRES_PORT"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	Host             string `env:"HOST"`
	Port             string `env:"PORT"`
}

var doOnce sync.Once
var Env ENV

func ParseEnv() *ENV {
	doOnce.Do(func() {
		err := godotenv.Load()

		if err != nil {
			log.Fatalf("Error while loading .env file: %s", err)
		}

		if err = env.Parse(&Env); err != nil {
			log.Fatalf("%+v\n", err)
		}

	})
	return &Env
}
