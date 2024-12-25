package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	PostgreSQL PostgreSQL
	JWT		   JWT
	App        Fiber
}

type Fiber struct {
	Host 	string
	Port 	string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type JWT struct {
	Secret	 string
}

func LoadConfigs() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Configs{
		PostgreSQL: PostgreSQL{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
		App: Fiber{
			Host: os.Getenv("APP_HOST"),
			Port: os.Getenv("APP_PORT"),
		},
		JWT: JWT{
			Secret:	os.Getenv("JWT_SECRET"),
		},
	}
}