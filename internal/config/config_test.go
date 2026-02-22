package config

import (
	"os"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	// 環境変数をクリアしてデフォルト値をテスト
	envVars := []string{
		"PORT", "SHUTDOWN_TIMEOUT", "CORS_ORIGINS",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE",
		"LOG_LEVEL", "LOG_FORMAT",
		"RATE_LIMIT_RPS", "RATE_LIMIT_BURST",
	}

	// 既存の環境変数を保存してクリア
	for _, key := range envVars {
		if val, ok := os.LookupEnv(key); ok {
			t.Setenv(key, val)
		}
		t.Setenv(key, "")
		os.Unsetenv(key) //nolint:errcheck // test cleanup
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Server defaults
	if cfg.Server.Port != "8080" {
		t.Errorf("Server.Port = %q, want %q", cfg.Server.Port, "8080")
	}
	if cfg.Server.ShutdownTimeout != 30 {
		t.Errorf("Server.ShutdownTimeout = %d, want %d", cfg.Server.ShutdownTimeout, 30)
	}
	if cfg.Server.CORSOrigins != "http://localhost:3000" {
		t.Errorf("Server.CORSOrigins = %q, want %q", cfg.Server.CORSOrigins, "http://localhost:3000")
	}

	// Database defaults
	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host = %q, want %q", cfg.Database.Host, "localhost")
	}
	if cfg.Database.Port != 55432 {
		t.Errorf("Database.Port = %d, want %d", cfg.Database.Port, 55432)
	}
	if cfg.Database.User != "postgres" {
		t.Errorf("Database.User = %q, want %q", cfg.Database.User, "postgres")
	}
	if cfg.Database.Password != "postgres" {
		t.Errorf("Database.Password = %q, want %q", cfg.Database.Password, "postgres")
	}
	if cfg.Database.DBName != "app_db" {
		t.Errorf("Database.DBName = %q, want %q", cfg.Database.DBName, "app_db")
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("Database.SSLMode = %q, want %q", cfg.Database.SSLMode, "disable")
	}

	// Log defaults
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "json")
	}

	// RateLimiter defaults
	if cfg.RateLimiter.RequestsPerSecond != 0 {
		t.Errorf("RateLimiter.RequestsPerSecond = %f, want %f", cfg.RateLimiter.RequestsPerSecond, 0.0)
	}
	if cfg.RateLimiter.BurstSize != 0 {
		t.Errorf("RateLimiter.BurstSize = %d, want %d", cfg.RateLimiter.BurstSize, 0)
	}
}

func TestLoad_EnvironmentVariableOverrides(t *testing.T) {
	// 環境変数を設定
	overrides := map[string]string{
		"PORT":             "9090",
		"SHUTDOWN_TIMEOUT": "60",
		"CORS_ORIGINS":     "https://example.com",
		"DB_HOST":          "db.example.com",
		"DB_PORT":          "5433",
		"DB_USER":          "myuser",
		"DB_PASSWORD":      "mypassword",
		"DB_NAME":          "mydb",
		"DB_SSLMODE":       "require",
		"LOG_LEVEL":        "debug",
		"LOG_FORMAT":       "text",
		"RATE_LIMIT_RPS":   "100.5",
		"RATE_LIMIT_BURST": "200",
	}

	for key, val := range overrides {
		t.Setenv(key, val)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Server overrides
	if cfg.Server.Port != "9090" {
		t.Errorf("Server.Port = %q, want %q", cfg.Server.Port, "9090")
	}
	if cfg.Server.ShutdownTimeout != 60 {
		t.Errorf("Server.ShutdownTimeout = %d, want %d", cfg.Server.ShutdownTimeout, 60)
	}
	if cfg.Server.CORSOrigins != "https://example.com" {
		t.Errorf("Server.CORSOrigins = %q, want %q", cfg.Server.CORSOrigins, "https://example.com")
	}

	// Database overrides
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("Database.Host = %q, want %q", cfg.Database.Host, "db.example.com")
	}
	if cfg.Database.Port != 5433 {
		t.Errorf("Database.Port = %d, want %d", cfg.Database.Port, 5433)
	}
	if cfg.Database.User != "myuser" {
		t.Errorf("Database.User = %q, want %q", cfg.Database.User, "myuser")
	}
	if cfg.Database.Password != "mypassword" {
		t.Errorf("Database.Password = %q, want %q", cfg.Database.Password, "mypassword")
	}
	if cfg.Database.DBName != "mydb" {
		t.Errorf("Database.DBName = %q, want %q", cfg.Database.DBName, "mydb")
	}
	if cfg.Database.SSLMode != "require" {
		t.Errorf("Database.SSLMode = %q, want %q", cfg.Database.SSLMode, "require")
	}

	// Log overrides
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "debug")
	}
	if cfg.Log.Format != "text" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "text")
	}

	// RateLimiter overrides
	if cfg.RateLimiter.RequestsPerSecond != 100.5 {
		t.Errorf("RateLimiter.RequestsPerSecond = %f, want %f", cfg.RateLimiter.RequestsPerSecond, 100.5)
	}
	if cfg.RateLimiter.BurstSize != 200 {
		t.Errorf("RateLimiter.BurstSize = %d, want %d", cfg.RateLimiter.BurstSize, 200)
	}
}
