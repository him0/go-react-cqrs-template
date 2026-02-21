package worker

import "time"

// Config ワーカーの設定
type Config struct {
	// PollInterval ポーリング間隔
	PollInterval time.Duration
	// BatchSize 1回のポーリングで取得するジョブ数
	BatchSize int
	// MaxConcurrency 同時実行するジョブの最大数
	MaxConcurrency int
	// ShutdownTimeout グレースフルシャットダウンのタイムアウト
	ShutdownTimeout time.Duration
}

// DefaultConfig デフォルト設定を返す
func DefaultConfig() Config {
	return Config{
		PollInterval:    5 * time.Second,
		BatchSize:       10,
		MaxConcurrency:  5,
		ShutdownTimeout: 30 * time.Second,
	}
}
