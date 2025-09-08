-- Add triggers to automatically update likes_count and dislikes_count in Posts and Comments tables

-- Triggers for Post_Reactions

-- Trigger for INSERT on Post_Reactions
CREATE TRIGGER update_post_reaction_counts_insert
AFTER INSERT ON Post_Reactions
BEGIN
    UPDATE Posts 
    SET likes_count = (
        SELECT COUNT(*) FROM Post_Reactions 
        WHERE post_id = NEW.post_id AND reaction_type = 'like'
    ),
    dislikes_count = (
        SELECT COUNT(*) FROM Post_Reactions 
        WHERE post_id = NEW.post_id AND reaction_type = 'dislike'
    )
    WHERE id = NEW.post_id;
END;

-- Trigger for UPDATE on Post_Reactions
CREATE TRIGGER update_post_reaction_counts_update
AFTER UPDATE ON Post_Reactions
BEGIN
    UPDATE Posts 
    SET likes_count = (
        SELECT COUNT(*) FROM Post_Reactions 
        WHERE post_id = NEW.post_id AND reaction_type = 'like'
    ),
    dislikes_count = (
        SELECT COUNT(*) FROM Post_Reactions 
        WHERE post_id = NEW.post_id AND reaction_type = 'dislike'
    )
    WHERE id = NEW.post_id;
END;

-- Trigger for DELETE on Post_Reactions
CREATE TRIGGER update_post_reaction_counts_delete
AFTER DELETE ON Post_Reactions
BEGIN
    UPDATE Posts 
    SET likes_count = (
        SELECT COUNT(*) FROM Post_Reactions 
        WHERE post_id = OLD.post_id AND reaction_type = 'like'
    ),
    dislikes_count = (
        SELECT COUNT(*) FROM Post_Reactions 
        WHERE post_id = OLD.post_id AND reaction_type = 'dislike'
    )
    WHERE id = OLD.post_id;
END;

-- Triggers for Comment_Reactions

-- Trigger for INSERT on Comment_Reactions
CREATE TRIGGER update_comment_reaction_counts_insert
AFTER INSERT ON Comment_Reactions
BEGIN
    UPDATE Comments 
    SET likes_count = (
        SELECT COUNT(*) FROM Comment_Reactions 
        WHERE comment_id = NEW.comment_id AND reaction_type = 'like'
    ),
    dislikes_count = (
        SELECT COUNT(*) FROM Comment_Reactions 
        WHERE comment_id = NEW.comment_id AND reaction_type = 'dislike'
    )
    WHERE id = NEW.comment_id;
END;

-- Trigger for UPDATE on Comment_Reactions
CREATE TRIGGER update_comment_reaction_counts_update
AFTER UPDATE ON Comment_Reactions
BEGIN
    UPDATE Comments 
    SET likes_count = (
        SELECT COUNT(*) FROM Comment_Reactions 
        WHERE comment_id = NEW.comment_id AND reaction_type = 'like'
    ),
    dislikes_count = (
        SELECT COUNT(*) FROM Comment_Reactions 
        WHERE comment_id = NEW.comment_id AND reaction_type = 'dislike'
    )
    WHERE id = NEW.comment_id;
END;

-- Trigger for DELETE on Comment_Reactions
CREATE TRIGGER update_comment_reaction_counts_delete
AFTER DELETE ON Comment_Reactions
BEGIN
    UPDATE Comments 
    SET likes_count = (
        SELECT COUNT(*) FROM Comment_Reactions 
        WHERE comment_id = OLD.comment_id AND reaction_type = 'like'
    ),
    dislikes_count = (
        SELECT COUNT(*) FROM Comment_Reactions 
        WHERE comment_id = OLD.comment_id AND reaction_type = 'dislike'
    )
    WHERE id = OLD.comment_id;
END;

-- Update existing counts to match current reactions
UPDATE Posts 
SET likes_count = (
    SELECT COUNT(*) FROM Post_Reactions 
    WHERE post_id = Posts.id AND reaction_type = 'like'
),
dislikes_count = (
    SELECT COUNT(*) FROM Post_Reactions 
    WHERE post_id = Posts.id AND reaction_type = 'dislike'
);

UPDATE Comments 
SET likes_count = (
    SELECT COUNT(*) FROM Comment_Reactions 
    WHERE comment_id = Comments.id AND reaction_type = 'like'
),
dislikes_count = (
    SELECT COUNT(*) FROM Comment_Reactions 
    WHERE comment_id = Comments.id AND reaction_type = 'dislike'
);
