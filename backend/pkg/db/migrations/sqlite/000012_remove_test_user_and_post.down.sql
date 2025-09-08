-- Restore test user ID 999 and associated data
-- Note: This is a basic restoration that recreates the user and post
-- Original relationships and detailed data cannot be fully restored

-- Recreate user 999
INSERT INTO Users (id, first_name, last_name, email, nickname, password, is_profile_public, created_at)
VALUES (999, 'Test', 'User', 'test999@test.com', 'testuser999', '$2a$10$hashedpassword', 1, datetime('now'));

-- Recreate post 999 (basic version)
INSERT INTO Posts (id, user_id, content, privacy, created_at)
VALUES (999, 999, 'Test post', 'public', datetime('now'));
