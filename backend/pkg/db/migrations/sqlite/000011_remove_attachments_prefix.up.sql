-- Remove 'attachments/' prefix from image paths in Posts table
UPDATE Posts 
SET image = SUBSTR(image, 13) 
WHERE image LIKE 'attachments/%';

-- Remove 'attachments/' prefix from image paths in Comments table  
UPDATE Comments 
SET image = SUBSTR(image, 13) 
WHERE image LIKE 'attachments/%';
