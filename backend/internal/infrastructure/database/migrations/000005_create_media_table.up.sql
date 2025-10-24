CREATE TABLE IF NOT EXISTS media (
    id VARCHAR(36) PRIMARY KEY,
    scene_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    url VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    file_size BIGINT,
    duration INT,
    width INT,
    height INT,
    format VARCHAR(50),
    metadata JSON,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE,
    INDEX idx_scene_id (scene_id),
    INDEX idx_type (type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
