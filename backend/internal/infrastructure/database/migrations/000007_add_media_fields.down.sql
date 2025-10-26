-- Remove added fields from aimotion_media table
DROP INDEX IF EXISTS idx_aimotion_media_completed_at;

ALTER TABLE aimotion_media
DROP COLUMN IF EXISTS completed_at,
DROP COLUMN IF EXISTS error_message,
DROP COLUMN IF EXISTS generation_id;

-- Revert column type changes
ALTER TABLE aimotion_media
ALTER COLUMN id TYPE VARCHAR(36),
ALTER COLUMN scene_id TYPE VARCHAR(36),
ALTER COLUMN url TYPE VARCHAR(500);
