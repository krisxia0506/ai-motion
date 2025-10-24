-- 场景表
CREATE TABLE aimotion_scene (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    chapter_id BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    sequence_num INT UNSIGNED NOT NULL COMMENT '场景序号',
    description TEXT COMMENT '场景描述',
    location VARCHAR(255) COMMENT '地点',
    time_of_day VARCHAR(50) COMMENT '时间段',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_novel_id (novel_id),
    INDEX idx_chapter_id (chapter_id),
    UNIQUE KEY uk_chapter_sequence (chapter_id, sequence_num)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景表';
