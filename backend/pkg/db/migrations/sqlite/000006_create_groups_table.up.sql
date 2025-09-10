-- Drop existing fragmented group tables
DROP TABLE IF EXISTS Group_Members;
DROP TABLE IF EXISTS Group_Posts;
DROP TABLE IF EXISTS Group_Events;
DROP TABLE IF EXISTS group_requests;
DROP TABLE IF EXISTS Group_Chat_Messages;
DROP TABLE IF EXISTS Group_Permissions;
DROP TABLE IF EXISTS Group_Post_Comments;
DROP TABLE IF EXISTS Group_Post_Reactions;
DROP TABLE IF EXISTS Group_Comment_Reactions;
DROP TABLE IF EXISTS Event_Responses;

-- Create unified Groups table
CREATE TABLE Groups (
    id TEXT PRIMARY KEY,
    public_id TEXT UNIQUE,
    type TEXT NOT NULL CHECK (type IN ('group', 'member', 'request', 'event', 'post', 'comment')),
    group_id TEXT,
    user_id INTEGER,
    title TEXT,
    content TEXT,
    role TEXT DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'pending', 'rejected', 'going', 'not_going')),
    privacy TEXT DEFAULT 'public' CHECK (privacy IN ('public', 'private')),
    image TEXT,
    data TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES Groups(id) ON DELETE CASCADE
);

-- Create performance indexes
CREATE INDEX idx_groups_public_id ON Groups(public_id);
CREATE INDEX idx_groups_type ON Groups(type);
CREATE INDEX idx_groups_group_id ON Groups(group_id);
CREATE INDEX idx_groups_user_id ON Groups(user_id);
CREATE INDEX idx_groups_status ON Groups(status);
CREATE INDEX idx_groups_type_group_id ON Groups(type, group_id);
CREATE INDEX idx_groups_type_user_id ON Groups(type, user_id);