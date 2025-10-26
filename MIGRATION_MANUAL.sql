-- ============================================================
-- Manual Migration: Add novel_id to aimotion_media table
-- ============================================================
-- Instructions:
-- 1. Go to your Supabase project dashboard
-- 2. Navigate to SQL Editor
-- 3. Copy and paste this entire script
-- 4. Run the script
-- ============================================================

-- Make scene_id nullable first (since manga workflow doesn't use scenes)
ALTER TABLE aimotion_media
ALTER COLUMN scene_id DROP NOT NULL;

-- Drop the foreign key constraint on scene_id since it's now optional
ALTER TABLE aimotion_media
DROP CONSTRAINT IF EXISTS aimotion_media_scene_id_fkey;

-- Add novel_id column to media table
ALTER TABLE aimotion_media
ADD COLUMN IF NOT EXISTS novel_id VARCHAR(36);

-- Add index for novel_id to improve query performance
CREATE INDEX IF NOT EXISTS idx_aimotion_media_novel_id ON aimotion_media(novel_id);

-- Add foreign key constraint to ensure referential integrity
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_aimotion_media_novel'
    ) THEN
        ALTER TABLE aimotion_media
        ADD CONSTRAINT fk_aimotion_media_novel
        FOREIGN KEY (novel_id) REFERENCES aimotion_novel(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Update comment
COMMENT ON COLUMN aimotion_media.novel_id IS '关联的小说ID，用于直接查询小说的所有媒体文件';

-- Verify the changes
SELECT
    column_name,
    data_type,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'aimotion_media'
AND column_name IN ('novel_id', 'scene_id')
ORDER BY ordinal_position;

-- Success message
DO $$
BEGIN
    RAISE NOTICE 'Migration completed successfully! The aimotion_media table now has a novel_id column.';
END $$;
