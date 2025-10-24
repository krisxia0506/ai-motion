-- 小说内容表 (大字段独立)
CREATE TABLE aimotion_novel_content (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    content LONGTEXT COMMENT '小说内容',
    
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    UNIQUE KEY uk_novel_id (novel_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小说内容表';
