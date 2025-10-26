-- Rollback: Remove novel_id from media table

-- Re-add foreign key constraint on scene_id
ALTER TABLE aimotion_media
ADD CONSTRAINT aimotion_media_scene_id_fkey
FOREIGN KEY (scene_id) REFERENCES aimotion_scene(id) ON DELETE CASCADE;

-- Make scene_id NOT NULL again
ALTER TABLE aimotion_media
ALTER COLUMN scene_id SET NOT NULL;

-- Drop foreign key constraint on novel_id
ALTER TABLE aimotion_media
DROP CONSTRAINT IF EXISTS fk_aimotion_media_novel;

-- Drop index
DROP INDEX IF EXISTS idx_aimotion_media_novel_id;

-- Remove novel_id column
ALTER TABLE aimotion_media
DROP COLUMN IF EXISTS novel_id;
