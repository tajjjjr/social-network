CREATE TABLE IF NOT EXISTS Comment_Reactions (
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    reaction_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
