-- name: EnqueueJob :exec
INSERT INTO jobs (id, job_type, payload, status, max_attempts, scheduled_at, created_at, updated_at)
VALUES ($1, $2, $3, 'pending', $4, $5, $6, $7);

-- name: FetchJobs :many
SELECT id, job_type, payload, status, attempts, max_attempts, last_error,
       scheduled_at, started_at, completed_at, created_at, updated_at
FROM jobs
WHERE status IN ('pending', 'retryable')
  AND scheduled_at <= NOW()
ORDER BY scheduled_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkJobProcessing :exec
UPDATE jobs
SET status = 'processing', started_at = NOW(), attempts = attempts + 1, updated_at = NOW()
WHERE id = $1;

-- name: MarkJobCompleted :exec
UPDATE jobs
SET status = 'completed', completed_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: MarkJobRetryable :exec
UPDATE jobs
SET status = 'retryable', last_error = $2, scheduled_at = $3, updated_at = NOW()
WHERE id = $1;

-- name: MarkJobDead :exec
UPDATE jobs
SET status = 'dead', last_error = $2, completed_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetJobByID :one
SELECT id, job_type, payload, status, attempts, max_attempts, last_error,
       scheduled_at, started_at, completed_at, created_at, updated_at
FROM jobs
WHERE id = $1;

-- name: ListJobsByStatus :many
SELECT id, job_type, payload, status, attempts, max_attempts, last_error,
       scheduled_at, started_at, completed_at, created_at, updated_at
FROM jobs
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountJobsByStatus :one
SELECT COUNT(*) FROM jobs WHERE status = $1;

-- name: DeleteCompletedJobsBefore :exec
DELETE FROM jobs
WHERE status IN ('completed', 'dead')
  AND completed_at < $1;
