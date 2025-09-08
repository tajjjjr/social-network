-- Add updated_at column to Posts table for tracking edits
ALTER TABLE Posts ADD COLUMN updated_at DATETIME;
