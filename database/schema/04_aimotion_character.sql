-- 角色表
CREATE TABLE aimotion_character (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    name VARCHAR(100) NOT NULL COMMENT '角色名称',
    appearance TEXT COMMENT '外貌描述',
    personality TEXT COMMENT '性格描述',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_novel_id (novel_id),
    INDEX idx_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';
