-- Create Group_Post_Comments table
CREATE TABLE IF NOT EXISTS Group_Post_Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    parent_comment_id INTEGER,
    content TEXT NOT NULL,
    image TEXT,
    like_count INTEGER DEFAULT 0,
    dislike_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_post_id) REFERENCES Group_Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES Group_Post_Comments(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_group_post_comments_group_post_id ON Group_Post_Comments(group_post_id);
CREATE INDEX IF NOT EXISTS idx_group_post_comments_user_id ON Group_Post_Comments(user_id);
