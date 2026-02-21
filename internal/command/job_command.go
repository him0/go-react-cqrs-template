package command

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/go-react-cqrs-template/internal/domain"
	"github.com/example/go-react-cqrs-template/internal/infrastructure"
	"github.com/example/go-react-cqrs-template/internal/infrastructure/dao"
)

// EnqueueJob ジョブをキューに追加（トランザクション内で使用）
func EnqueueJob(ctx context.Context, tx infrastructure.DBTX, job *domain.Job) error {
	queries := dao.New(tx)
	err := queries.EnqueueJob(ctx, dao.EnqueueJobParams{
		ID:          job.ID,
		JobType:     job.JobType,
		Payload:     job.Payload,
		MaxAttempts: int32(job.MaxAttempts),
		ScheduledAt: job.ScheduledAt,
		CreatedAt:   job.CreatedAt,
		UpdatedAt:   job.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}
	return nil
}

// FetchAndLockJobs ジョブを取得しロック（トランザクション内で使用）
func FetchAndLockJobs(ctx context.Context, tx infrastructure.DBTX, limit int) ([]*domain.Job, error) {
	queries := dao.New(tx)
	rows, err := queries.FetchJobs(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}

	jobs := make([]*domain.Job, 0, len(rows))
	for _, row := range rows {
		jobs = append(jobs, toDomainJob(row))
	}
	return jobs, nil
}

// MarkJobProcessing ジョブを処理中に変更（トランザクション内で使用）
func MarkJobProcessing(ctx context.Context, tx infrastructure.DBTX, jobID string) error {
	queries := dao.New(tx)
	return queries.MarkJobProcessing(ctx, jobID)
}

// MarkJobCompleted ジョブを完了に変更（トランザクション内で使用）
func MarkJobCompleted(ctx context.Context, tx infrastructure.DBTX, jobID string) error {
	queries := dao.New(tx)
	return queries.MarkJobCompleted(ctx, jobID)
}

// MarkJobRetryable ジョブをリトライ可能に変更（トランザクション内で使用）
func MarkJobRetryable(ctx context.Context, tx infrastructure.DBTX, jobID string, lastError string, nextScheduledAt time.Time) error {
	queries := dao.New(tx)
	return queries.MarkJobRetryable(ctx, dao.MarkJobRetryableParams{
		ID:          jobID,
		LastError:   sql.NullString{String: lastError, Valid: true},
		ScheduledAt: nextScheduledAt,
	})
}

// MarkJobDead ジョブをデッドに変更（トランザクション内で使用）
func MarkJobDead(ctx context.Context, tx infrastructure.DBTX, jobID string, lastError string) error {
	queries := dao.New(tx)
	return queries.MarkJobDead(ctx, dao.MarkJobDeadParams{
		ID:        jobID,
		LastError: sql.NullString{String: lastError, Valid: true},
	})
}

// toDomainJob dao.Jobをdomain.Jobに変換
func toDomainJob(j dao.Job) *domain.Job {
	job := &domain.Job{
		ID:          j.ID,
		JobType:     j.JobType,
		Payload:     json.RawMessage(j.Payload),
		Status:      domain.JobStatus(j.Status),
		Attempts:    int(j.Attempts),
		MaxAttempts: int(j.MaxAttempts),
		ScheduledAt: j.ScheduledAt,
		CreatedAt:   j.CreatedAt,
		UpdatedAt:   j.UpdatedAt,
	}
	if j.LastError.Valid {
		job.LastError = j.LastError.String
	}
	if j.StartedAt.Valid {
		job.StartedAt = &j.StartedAt.Time
	}
	if j.CompletedAt.Valid {
		job.CompletedAt = &j.CompletedAt.Time
	}
	return job
}
