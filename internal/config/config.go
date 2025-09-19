package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
	DBDsn string
}

func LoadConfig() *Config{
	_ = godotenv.Load()

	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "3000")
	viper.SetDefault("DB_DSN","postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" )
	
	cfg := &Config{
		AppPort: viper.GetString("APP_PORT"),
		DBDsn: viper.GetString("DB_DSN"),
	}

	log.Printf("Loaded config: port=%s db=%s", cfg.AppPort, cfg.DBDsn)
	return cfg
}