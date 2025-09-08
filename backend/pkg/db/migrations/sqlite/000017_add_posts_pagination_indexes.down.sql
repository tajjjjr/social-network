-- Remove pagination indexes
DROP INDEX IF EXISTS idx_posts_created_at_desc;
DROP INDEX IF EXISTS idx_posts_privacy_created_at;
DROP INDEX IF EXISTS idx_posts_user_privacy_created_at;