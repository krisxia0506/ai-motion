-- 角色图片表
CREATE TABLE aimotion_character_image (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    character_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    image_url VARCHAR(512) NOT NULL COMMENT '图片URL',
    image_type TINYINT UNSIGNED DEFAULT 0 COMMENT '图片类型:0-参考图,1-场景图',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_character_id (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色图片表';
