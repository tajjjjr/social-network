-- Create Groups table
CREATE TABLE
    IF NOT EXISTS Groups (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        avatar TEXT,
        title TEXT NOT NULL,
        description TEXT,
        privacy TEXT NOT NULL DEFAULT 'public',
        creator_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (creator_id) REFERENCES Users (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_groups_creator_id ON Groups (creator_id);