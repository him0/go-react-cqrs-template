-- User logs table
CREATE TABLE IF NOT EXISTS user_logs (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index for user_id lookup
CREATE INDEX IF NOT EXISTS idx_user_logs_user_id ON user_logs(user_id);

-- Index for action lookup
CREATE INDEX IF NOT EXISTS idx_user_logs_action ON user_logs(action);

-- Index for created_at for sorting
CREATE INDEX IF NOT EXISTS idx_user_logs_created_at ON user_logs(created_at DESC);
