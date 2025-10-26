-- Add novel_id to media table to establish direct relationship between media and novel
ALTER TABLE aimotion_media
ADD COLUMN novel_id VARCHAR(36);

-- Add index for novel_id to improve query performance
CREATE INDEX IF NOT EXISTS idx_aimotion_media_novel_id ON aimotion_media(novel_id);

-- Add foreign key constraint to ensure referential integrity
ALTER TABLE aimotion_media
ADD CONSTRAINT fk_aimotion_media_novel
FOREIGN KEY (novel_id) REFERENCES aimotion_novel(id) ON DELETE CASCADE;

-- Update comment
COMMENT ON COLUMN aimotion_media.novel_id IS '关联的小说ID，用于直接查询小说的所有媒体文件';

-- Make scene_id nullable since manga workflow doesn't use scenes
ALTER TABLE aimotion_media
ALTER COLUMN scene_id DROP NOT NULL;

-- Drop the foreign key constraint on scene_id since it's now optional
ALTER TABLE aimotion_media
DROP CONSTRAINT IF EXISTS aimotion_media_scene_id_fkey;
