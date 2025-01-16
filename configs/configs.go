package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct {
	PostgreSQL PostgreSQL
	JWT        JWT
	App        Fiber
	Supabase   Supabase
	Mail       Mail
}

type Fiber struct {
	Host string
	Port string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type JWT struct {
	Secret string
}

type Supabase struct {
	URL    string
	Key    string
	Bucket string
}

type Mail struct {
	Host   string
	Port   string
	Sender string
	Key    string
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
			SSLMode:  os.Getenv("SSL_Mode"),
		},
		App: Fiber{
			Host: os.Getenv("APP_HOST"),
			Port: os.Getenv("APP_PORT"),
		},
		JWT: JWT{
			Secret: os.Getenv("JWT_SECRET"),
		},
		Supabase: Supabase{
			URL:    os.Getenv("SUPABASE_URL"),
			Key:    os.Getenv("SUPABASE_KEY"),
			Bucket: os.Getenv("BUCKET_NAME"),
		},
		Mail: Mail{
			Host:   os.Getenv("EMAIL_HOST"),
			Port:   os.Getenv("EMAIL_PORT"),
			Sender: os.Getenv("EMAIL_USER"),
			Key:    os.Getenv("EMAIL_PASS"),
		},
	}
}
