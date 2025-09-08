-- Add indexes for efficient pagination queries
CREATE INDEX IF NOT EXISTS idx_posts_created_at_desc ON Posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_privacy_created_at ON Posts(privacy, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_user_privacy_created_at ON Posts(user_id, privacy, created_at DESC);