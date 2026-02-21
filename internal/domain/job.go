package domain

import (
	"crypto/rand"
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
)

// JobStatus ジョブの状態
type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusRetryable  JobStatus = "retryable"
	JobStatusDead       JobStatus = "dead"
)

// Job ジョブのドメインモデル
type Job struct {
	ID          string
	JobType     string
	Payload     json.RawMessage
	Status      JobStatus
	Attempts    int
	MaxAttempts int
	LastError   string
	ScheduledAt time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewJob 新しいジョブを作成
func NewJob(jobType string, payload json.RawMessage, maxAttempts int) *Job {
	now := time.Now()
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	return &Job{
		ID:          ulid.MustNew(ulid.Timestamp(now), rand.Reader).String(),
		JobType:     jobType,
		Payload:     payload,
		Status:      JobStatusPending,
		MaxAttempts: maxAttempts,
		ScheduledAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewScheduledJob スケジュール付きジョブを作成
func NewScheduledJob(jobType string, payload json.RawMessage, maxAttempts int, scheduledAt time.Time) *Job {
	job := NewJob(jobType, payload, maxAttempts)
	job.ScheduledAt = scheduledAt
	return job
}

// CanRetry リトライ可能かどうか判定
func (j *Job) CanRetry() bool {
	return j.Attempts < j.MaxAttempts
}
