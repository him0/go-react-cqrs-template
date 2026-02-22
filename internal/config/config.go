package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Config はアプリケーション全体の設定を保持する構造体
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Log         LogConfig
	RateLimiter RateLimiterConfig
}

// ServerConfig はHTTPサーバーの設定
type ServerConfig struct {
	Port            string `envconfig:"PORT" default:"8080"`
	ShutdownTimeout int    `envconfig:"SHUTDOWN_TIMEOUT" default:"30"`
	CORSOrigins     string `envconfig:"CORS_ORIGINS" default:"http://localhost:3000"`
}

// DatabaseConfig はデータベース接続の設定
type DatabaseConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     int    `envconfig:"DB_PORT" default:"55432"`
	User     string `envconfig:"DB_USER" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName   string `envconfig:"DB_NAME" default:"app_db"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
}

// LogConfig はロギングの設定
type LogConfig struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
}

// RateLimiterConfig はレートリミッターの設定
type RateLimiterConfig struct {
	RequestsPerSecond float64 `envconfig:"RATE_LIMIT_RPS" default:"0"`
	BurstSize         int     `envconfig:"RATE_LIMIT_BURST" default:"0"`
}

// Load は環境変数からConfigを読み込む
func Load() (*Config, error) {
	var cfg Config
	if err := processEnvConfig(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
}

// processEnvConfig は構造体のenvconfigタグを読み取り、環境変数から値を設定する
func processEnvConfig(cfg any) error {
	return processStruct(reflect.ValueOf(cfg).Elem())
}

// processStruct は構造体の各フィールドを再帰的に処理する
func processStruct(v reflect.Value) error {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// ネストされた構造体を再帰処理
		if field.Type.Kind() == reflect.Struct {
			if err := processStruct(fieldVal); err != nil {
				return err
			}
			continue
		}

		envKey := field.Tag.Get("envconfig")
		if envKey == "" {
			continue
		}

		defaultVal := field.Tag.Get("default")
		envVal := os.Getenv(envKey)

		// 環境変数が設定されていなければデフォルト値を使用
		val := envVal
		if val == "" {
			val = defaultVal
		}

		if val == "" {
			continue
		}

		// フィールドの型に応じて値を設定
		switch field.Type.Kind() {
		case reflect.String:
			fieldVal.SetString(val)
		case reflect.Int:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("failed to parse %s=%q as int: %w", envKey, val, err)
			}
			fieldVal.SetInt(int64(intVal))
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("failed to parse %s=%q as float64: %w", envKey, val, err)
			}
			fieldVal.SetFloat(floatVal)
		default:
			return fmt.Errorf("unsupported field type %s for %s", field.Type.Kind(), field.Name)
		}
	}

	return nil
}
