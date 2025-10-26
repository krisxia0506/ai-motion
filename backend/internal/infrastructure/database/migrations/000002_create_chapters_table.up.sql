CREATE TABLE IF NOT EXISTS aimotion_chapter (
    id VARCHAR(36) PRIMARY KEY,
    novel_id VARCHAR(36) NOT NULL,
    chapter_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    word_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (novel_id) REFERENCES aimotion_novel(id) ON DELETE CASCADE,
    UNIQUE (novel_id, chapter_number)
);

CREATE INDEX IF NOT EXISTS idx_aimotion_chapter_novel_id ON aimotion_chapter(novel_id);
CREATE INDEX IF NOT EXISTS idx_aimotion_chapter_chapter_number ON aimotion_chapter(chapter_number);

-- Trigger to automatically update updated_at
CREATE OR REPLACE FUNCTION update_aimotion_chapter_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_aimotion_chapter_updated_at
    BEFORE UPDATE ON aimotion_chapter
    FOR EACH ROW
    EXECUTE FUNCTION update_aimotion_chapter_updated_at();
