-- Remove triggers for automatic reaction count updates

-- Drop Post_Reactions triggers
DROP TRIGGER IF EXISTS update_post_reaction_counts_insert;
DROP TRIGGER IF EXISTS update_post_reaction_counts_update;
DROP TRIGGER IF EXISTS update_post_reaction_counts_delete;

-- Drop Comment_Reactions triggers
DROP TRIGGER IF EXISTS update_comment_reaction_counts_insert;
DROP TRIGGER IF EXISTS update_comment_reaction_counts_update;
DROP TRIGGER IF EXISTS update_comment_reaction_counts_delete;

-- Reset counts to 0
UPDATE Posts SET likes_count = 0, dislikes_count = 0;
UPDATE Comments SET likes_count = 0, dislikes_count = 0;
