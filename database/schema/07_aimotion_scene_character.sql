-- 场景角色关联表
CREATE TABLE aimotion_scene_character (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    scene_id BIGINT UNSIGNED NOT NULL COMMENT '场景ID',
    character_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    UNIQUE KEY uk_scene_character (scene_id, character_id),
    INDEX idx_character_id (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景角色关联表';
