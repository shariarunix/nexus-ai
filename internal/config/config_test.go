package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Backup original env
	originalEnv := make(map[string]string)
	envVars := []string{"SUPABASE_URL", "SUPABASE_KEY", "GEMINI_API_KEY", "PORT"}
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
		os.Unsetenv(key)
	}
	defer func() {
		for key, val := range originalEnv {
			os.Setenv(key, val)
		}
	}()

	t.Run("Success loading minimal config", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://user:pass@host/db")
		os.Setenv("SUPABASE_URL", "https://example.supabase.co")
		os.Setenv("GEMINI_API_KEY", "gemini-key")

		cfg, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, "postgres://user:pass@host/db", cfg.DatabaseURL)
		assert.Equal(t, "https://example.supabase.co", cfg.SupabaseURL)
		assert.Equal(t, "gemini-key", cfg.GeminiAPIKey)
		assert.Equal(t, "8080", cfg.Port) // Default port
	})

	t.Run("Success loading with custom port", func(t *testing.T) {
		os.Setenv("DATABASE_URL", "postgres://user:pass@host/db")
		os.Setenv("SUPABASE_URL", "https://example.supabase.co")
		os.Setenv("GEMINI_API_KEY", "gemini-key")
		os.Setenv("PORT", "9090")

		cfg, err := Load()
		assert.NoError(t, err)
		assert.Equal(t, "9090", cfg.Port)
	})

	t.Run("Failure missing required env", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")

		_, err := Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "DATABASE_URL is required")
	})
}
