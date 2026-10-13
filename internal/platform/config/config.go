package config 

import (
	"fmt"
	"os"
)



// Config holds everything the service needs from it environment 
type Config struct {
	HTTPAddr string 
	DatabaseURL string 
}


//Load reads config from environment variables, with defaults for local dev.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		return Config{} , fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}


func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
