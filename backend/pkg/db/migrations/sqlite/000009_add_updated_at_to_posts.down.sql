-- Remove updated_at column from Posts table
-- Note: SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
CREATE TABLE Posts_backup AS SELECT id, user_id, content, image, privacy, created_at FROM Posts;
DROP TABLE Posts;
CREATE TABLE Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    privacy TEXT NOT NULL DEFAULT 'public',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);
INSERT INTO Posts SELECT * FROM Posts_backup;
DROP TABLE Posts_backup;
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON Posts(user_id);
