-- +migrate Up
ALTER TABLE users
 ADD COLUMN avatar_imagekit_file_id TEXT,
 ADD COLUMN avatar_imagekit_file_path TEXT,
 ADD CONSTRAINT users_managed_avatar_complete CHECK (
   (avatar_imagekit_file_id IS NULL AND avatar_imagekit_file_path IS NULL)
   OR (avatar_imagekit_file_id IS NOT NULL AND length(btrim(avatar_imagekit_file_id)) > 0
       AND avatar_imagekit_file_path IS NOT NULL AND length(btrim(avatar_imagekit_file_path)) > 0
       AND avatar_url IS NOT NULL AND length(btrim(avatar_url)) > 0)
 );
ALTER TABLE visitors
 ADD COLUMN photo_imagekit_file_id TEXT,
 ADD COLUMN photo_imagekit_file_path TEXT,
 ADD CONSTRAINT visitors_managed_photo_complete CHECK (
   (photo_imagekit_file_id IS NULL AND photo_imagekit_file_path IS NULL)
   OR (photo_imagekit_file_id IS NOT NULL AND length(btrim(photo_imagekit_file_id)) > 0
       AND photo_imagekit_file_path IS NOT NULL AND length(btrim(photo_imagekit_file_path)) > 0
       AND photo_url IS NOT NULL AND length(btrim(photo_url)) > 0)
 );

-- +migrate Down
ALTER TABLE visitors DROP CONSTRAINT visitors_managed_photo_complete,
 DROP COLUMN photo_imagekit_file_path, DROP COLUMN photo_imagekit_file_id;
ALTER TABLE users DROP CONSTRAINT users_managed_avatar_complete,
 DROP COLUMN avatar_imagekit_file_path, DROP COLUMN avatar_imagekit_file_id;
