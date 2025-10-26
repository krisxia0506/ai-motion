CREATE TABLE IF NOT EXISTS aimotion_media (
    id VARCHAR(36) PRIMARY KEY,
    scene_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    url VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    file_size BIGINT,
    duration INTEGER,
    width INTEGER,
    height INTEGER,
    format VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (scene_id) REFERENCES aimotion_scene(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_aimotion_media_scene_id ON aimotion_media(scene_id);
CREATE INDEX IF NOT EXISTS idx_aimotion_media_type ON aimotion_media(type);
CREATE INDEX IF NOT EXISTS idx_aimotion_media_status ON aimotion_media(status);

-- Trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_aimotion_media_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_aimotion_media_updated_at
    BEFORE UPDATE ON aimotion_media
    FOR EACH ROW
    EXECUTE FUNCTION update_aimotion_media_updated_at();
