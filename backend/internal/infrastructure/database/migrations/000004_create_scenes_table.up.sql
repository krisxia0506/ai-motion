CREATE TABLE IF NOT EXISTS scenes (
    id VARCHAR(36) PRIMARY KEY,
    chapter_id VARCHAR(36) NOT NULL,
    scene_number INT NOT NULL,
    description TEXT NOT NULL,
    dialogue TEXT,
    location VARCHAR(255),
    time_of_day VARCHAR(50),
    characters JSON,
    prompt TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE,
    INDEX idx_chapter_id (chapter_id),
    INDEX idx_scene_number (scene_number),
    UNIQUE KEY uk_chapter_scene (chapter_id, scene_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
