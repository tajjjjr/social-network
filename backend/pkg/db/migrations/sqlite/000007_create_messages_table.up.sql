-- Create Messages table
CREATE TABLE IF NOT EXISTS Messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER,
    group_id INTEGER,
    content TEXT,
    is_emoji BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON Messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON Messages(receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_group_id ON Messages(group_id);
