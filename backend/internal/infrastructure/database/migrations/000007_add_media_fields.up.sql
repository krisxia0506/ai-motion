-- Add missing fields to aimotion_media table
ALTER TABLE aimotion_media
ADD COLUMN IF NOT EXISTS generation_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS error_message TEXT,
ADD COLUMN IF NOT EXISTS completed_at TIMESTAMP;

-- Drop foreign key constraint (scene feature not needed)
ALTER TABLE aimotion_media
DROP CONSTRAINT IF EXISTS aimotion_media_scene_id_fkey;

-- Extend id and scene_id column lengths to support longer identifiers
ALTER TABLE aimotion_media
ALTER COLUMN id TYPE VARCHAR(100),
ALTER COLUMN scene_id TYPE VARCHAR(100);

-- Extend url column to support longer URLs (e.g., long cloud storage URLs or data URLs)
ALTER TABLE aimotion_media
ALTER COLUMN url TYPE TEXT;

-- Create index on completed_at for better query performance
CREATE INDEX IF NOT EXISTS idx_aimotion_media_completed_at ON aimotion_media(completed_at);
