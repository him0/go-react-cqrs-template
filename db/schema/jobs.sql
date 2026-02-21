-- ジョブキューテーブル
CREATE TABLE IF NOT EXISTS jobs (
    id VARCHAR(26) PRIMARY KEY,
    job_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    last_error TEXT,
    scheduled_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ポーリング用部分インデックス（pending/retryable のみ）
CREATE INDEX IF NOT EXISTS idx_jobs_pollable
    ON jobs(scheduled_at ASC)
    WHERE status IN ('pending', 'retryable');

-- ステータス別クエリ用
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);

-- ジョブタイプ別クエリ用
CREATE INDEX IF NOT EXISTS idx_jobs_job_type ON jobs(job_type);
