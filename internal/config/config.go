package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
	DBDsn   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_DSN", "postgres://user:1234@localhost:5432/mydb?sslmode=disable")


	cfg := &Config{
		AppPort: viper.GetString("APP_PORT"),
		DBDsn:   viper.GetString("DB_DSN"),
	}

	log.Printf("Loaded config: port=%s db=%s", cfg.AppPort, cfg.DBDsn)

	// Optional: sanity check
	if cfg.DBDsn == "" {
		log.Fatal("DB_DSN is empty. Please set it in .env or Docker environment.")
		os.Exit(1)
	}

	return cfg
}
