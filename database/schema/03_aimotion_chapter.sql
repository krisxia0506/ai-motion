-- 章节表
CREATE TABLE aimotion_chapter (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    title VARCHAR(255) COMMENT '章节标题',
    sequence_num INT UNSIGNED NOT NULL COMMENT '章节序号',
    content TEXT COMMENT '章节内容',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_novel_id (novel_id),
    UNIQUE KEY uk_novel_sequence (novel_id, sequence_num)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节表';
