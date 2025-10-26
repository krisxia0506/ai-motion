CREATE TABLE IF NOT EXISTS aimotion_novel (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    word_count INTEGER NOT NULL DEFAULT 0,
    chapter_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_aimotion_novel_status ON aimotion_novel(status);
CREATE INDEX IF NOT EXISTS idx_aimotion_novel_created_at ON aimotion_novel(created_at);

-- Trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_aimotion_novel_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_aimotion_novel_updated_at
    BEFORE UPDATE ON aimotion_novel
    FOR EACH ROW
    EXECUTE FUNCTION update_aimotion_novel_updated_at();
