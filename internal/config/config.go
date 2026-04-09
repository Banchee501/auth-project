package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	accessTTL, _ := time.ParseDuration(os.Getenv("ACCESS_TTL"))
	refreshTTL, _ := time.ParseDuration(os.Getenv("REFRESH_TTL"))

	return &Config{
		DBUrl:      os.Getenv("DB_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AccessTTL:  accessTTL,
		RefreshTTL: refreshTTL,
	}
}
