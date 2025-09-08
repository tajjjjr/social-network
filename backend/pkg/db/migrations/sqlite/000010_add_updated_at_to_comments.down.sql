-- Remove updated_at column from Comments table
-- Note: SQLite doesn't support DROP COLUMN directly, so we need to recreate the table

-- Create temporary table with original schema
CREATE TABLE Comments_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Copy data from original table (excluding updated_at)
INSERT INTO Comments_temp (id, post_id, user_id, content, image, created_at)
SELECT id, post_id, user_id, content, image, created_at FROM Comments;

-- Drop original table
DROP TABLE Comments;

-- Rename temp table to original name
ALTER TABLE Comments_temp RENAME TO Comments;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON Comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON Comments(user_id);