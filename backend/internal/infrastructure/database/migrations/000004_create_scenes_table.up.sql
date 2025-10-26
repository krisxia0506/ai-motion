CREATE TABLE IF NOT EXISTS aimotion_scene (
    id VARCHAR(36) PRIMARY KEY,
    chapter_id VARCHAR(36) NOT NULL,
    scene_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    dialogue TEXT,
    location VARCHAR(255),
    time_of_day VARCHAR(50),
    characters JSONB,
    prompt TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chapter_id) REFERENCES aimotion_chapter(id) ON DELETE CASCADE,
    UNIQUE (chapter_id, scene_number)
);

CREATE INDEX IF NOT EXISTS idx_aimotion_scene_chapter_id ON aimotion_scene(chapter_id);
CREATE INDEX IF NOT EXISTS idx_aimotion_scene_scene_number ON aimotion_scene(scene_number);

-- Trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_aimotion_scene_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_aimotion_scene_updated_at
    BEFORE UPDATE ON aimotion_scene
    FOR EACH ROW
    EXECUTE FUNCTION update_aimotion_scene_updated_at();
