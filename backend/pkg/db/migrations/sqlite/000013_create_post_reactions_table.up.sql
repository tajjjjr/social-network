CREATE TABLE IF NOT EXISTS Post_Reactions (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    reaction_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (post_id) REFERENCES Posts(id)
);
