-- Add 'attachments/' prefix back to image paths in Posts table
UPDATE Posts 
SET image = 'attachments/' || image 
WHERE image IS NOT NULL AND image != '' AND image NOT LIKE 'attachments/%';

-- Add 'attachments/' prefix back to image paths in Comments table
UPDATE Comments 
SET image = 'attachments/' || image 
WHERE image IS NOT NULL AND image != '' AND image NOT LIKE 'attachments/%';
