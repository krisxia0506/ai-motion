CREATE TABLE IF NOT EXISTS aimotion_character (
    id VARCHAR(36) PRIMARY KEY,
    novel_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    appearance TEXT,
    personality TEXT,
    description TEXT,
    reference_image_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (novel_id) REFERENCES aimotion_novel(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_aimotion_character_novel_id ON aimotion_character(novel_id);
CREATE INDEX IF NOT EXISTS idx_aimotion_character_role ON aimotion_character(role);
CREATE INDEX IF NOT EXISTS idx_aimotion_character_name ON aimotion_character(name);

-- Trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_aimotion_character_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_aimotion_character_updated_at
    BEFORE UPDATE ON aimotion_character
    FOR EACH ROW
    EXECUTE FUNCTION update_aimotion_character_updated_at();
