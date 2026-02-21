package worker

import (
	"context"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/example/go-react-cqrs-template/internal/command"
	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/infrastructure"
)

// Worker 非同期ジョブワーカー
type Worker struct {
	txManager *infrastructure.TransactionManager
	registry  *Registry
	config    Config
	logger    *slog.Logger
}

// NewWorker Workerのコンストラクタ
func NewWorker(
	txManager *infrastructure.TransactionManager,
	registry *Registry,
	config Config,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		txManager: txManager,
		registry:  registry,
		config:    config,
		logger:    logger,
	}
}

// Run ワーカーを起動（ctx がキャンセルされるまで実行し続ける）
func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("worker started",
		slog.Duration("poll_interval", w.config.PollInterval),
		slog.Int("batch_size", w.config.BatchSize),
		slog.Int("max_concurrency", w.config.MaxConcurrency),
	)

	sem := make(chan struct{}, w.config.MaxConcurrency)
	var wg sync.WaitGroup

	ticker := time.NewTicker(w.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker shutting down, waiting for in-flight jobs...")
			wg.Wait()
			w.logger.Info("worker stopped")
			return nil
		case <-ticker.C:
			w.poll(ctx, sem, &wg)
		}
	}
}

func (w *Worker) poll(ctx context.Context, sem chan struct{}, wg *sync.WaitGroup) {
	err := w.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
		jobs, err := command.FetchAndLockJobs(ctx, tx, w.config.BatchSize)
		if err != nil {
			return err
		}

		for _, job := range jobs {
			if err := command.MarkJobProcessing(ctx, tx, job.ID); err != nil {
				w.logger.Error("failed to mark job processing",
					slog.String("job_id", job.ID),
					slog.String("error", err.Error()),
				)
				continue
			}

			job := job
			wg.Add(1)
			sem <- struct{}{} // セマフォスロット取得

			go func() {
				defer wg.Done()
				defer func() { <-sem }() // セマフォスロット解放

				w.processJob(ctx, job)
			}()
		}

		return nil
	})

	if err != nil {
		if ctx.Err() != nil {
			return // コンテキストキャンセル時はエラーではない
		}
		w.logger.Error("failed to poll jobs",
			slog.String("error", err.Error()),
		)
	}
}

func (w *Worker) processJob(ctx context.Context, job *domain.Job) {
	jobLogger := w.logger.With(
		slog.String("job_id", job.ID),
		slog.String("job_type", job.JobType),
		slog.Int("attempt", job.Attempts),
	)

	jobLogger.Info("processing job")
	startTime := time.Now()

	handler, err := w.registry.Get(job.JobType)
	if err != nil {
		jobLogger.Error("no handler for job type", slog.String("error", err.Error()))
		_ = w.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
			return command.MarkJobDead(ctx, tx, job.ID, err.Error())
		})
		return
	}

	if err := handler.Handle(ctx, job.Payload); err != nil {
		duration := time.Since(startTime)
		jobLogger.Error("job failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", duration),
		)

		_ = w.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
			if job.CanRetry() {
				backoff := calculateBackoff(job.Attempts)
				nextScheduledAt := time.Now().Add(backoff)
				jobLogger.Info("scheduling retry",
					slog.Duration("backoff", backoff),
					slog.Time("next_scheduled_at", nextScheduledAt),
				)
				return command.MarkJobRetryable(ctx, tx, job.ID, err.Error(), nextScheduledAt)
			}
			jobLogger.Error("job moved to dead letter (max attempts reached)")
			return command.MarkJobDead(ctx, tx, job.ID, err.Error())
		})
		return
	}

	duration := time.Since(startTime)
	jobLogger.Info("job completed", slog.Duration("duration", duration))

	_ = w.txManager.RunInTransaction(ctx, func(ctx context.Context, tx infrastructure.DBTX) error {
		return command.MarkJobCompleted(ctx, tx, job.ID)
	})
}

// calculateBackoff 指数バックオフを計算（base 5s, max 5min）
func calculateBackoff(attempt int) time.Duration {
	base := 5.0 // seconds
	backoff := base * math.Pow(2, float64(attempt-1))
	maxBackoff := 300.0 // 5 minutes
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	return time.Duration(backoff * float64(time.Second))
}
