-- Create Group_Permissions table for admin permissions
CREATE TABLE IF NOT EXISTS Group_Permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    permission_type TEXT NOT NULL, -- 'admin', 'moderator'
    granted_by INTEGER NOT NULL,
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES Users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id, permission_type)
);

CREATE INDEX IF NOT EXISTS idx_group_permissions_group_id ON Group_Permissions(group_id);
CREATE INDEX IF NOT EXISTS idx_group_permissions_user_id ON Group_Permissions(user_id);