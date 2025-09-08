-- Add updated_at column to Comments table for tracking edits
ALTER TABLE Comments ADD COLUMN updated_at DATETIME;

-- Update existing comments to have updated_at = created_at
UPDATE Comments SET updated_at = created_at WHERE updated_at IS NULL;