package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Environment            string
	Port                   string
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string
	SupabaseJWTSecret      string
	DatabaseURL            string
	GeminiAPIKey           string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (optional)
	_ = godotenv.Load("../../.env")

	cfg := &Config{
		Environment:            getEnv("ENV", "development"),
		Port:                   getEnv("PORT", "8080"),
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:        os.Getenv("SUPABASE_ANON_KEY"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SupabaseJWTSecret:      os.Getenv("SUPABASE_JWT_SECRET"),
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.GeminiAPIKey == "" {
		return nil, errors.New("GEMINI_API_KEY is required")
	}
	if cfg.SupabaseURL == "" {
		return nil, errors.New("SUPABASE_URL is required")
	}
	if cfg.SupabaseAnonKey == "" {
		return nil, errors.New("SUPABASE_ANON_KEY is required")
	}
	if cfg.SupabaseServiceRoleKey == "" {
		return nil, errors.New("SUPABASE_SERVICE_ROLE_KEY is required")
	}
	if cfg.SupabaseJWTSecret == "" {
		return nil, errors.New("SUPABASE_JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
