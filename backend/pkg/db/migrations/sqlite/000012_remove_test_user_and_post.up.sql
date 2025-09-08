-- Remove test user ID 999 and associated data

-- Delete comments made by user 999
DELETE FROM Comments WHERE user_id = 999;

-- Delete comments on posts by user 999
DELETE FROM Comments WHERE post_id IN (SELECT id FROM Posts WHERE user_id = 999);

-- Delete post visibility records for posts by user 999
DELETE FROM Post_visibility WHERE post_id IN (SELECT id FROM Posts WHERE user_id = 999);

-- Delete post with ID 999 specifically
DELETE FROM Posts WHERE id = 999;

-- Delete all posts by user 999
DELETE FROM Posts WHERE user_id = 999;

-- Delete follower relationships involving user 999
DELETE FROM Followers WHERE follower_id = 999 OR followee_id = 999;

-- Delete sessions for user 999
DELETE FROM Sessions WHERE user_id = 999;

-- Delete notifications for user 999
-- DELETE FROM Notifications WHERE user_id = 999 OR from_user_id = 999;

-- Delete messages involving user 999
DELETE FROM Messages WHERE sender_id = 999 OR receiver_id = 999;

-- Delete group memberships for user 999 (unified table)
DELETE FROM Groups WHERE type = 'member' AND user_id = 999;

-- Delete event responses by user 999 (unified table)
DELETE FROM Groups WHERE type = 'event_response' AND user_id = 999;

-- Finally, delete user 999
DELETE FROM Users WHERE id = 999;
