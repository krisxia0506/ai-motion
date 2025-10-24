-- 媒体表
CREATE TABLE aimotion_media (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    type TINYINT UNSIGNED NOT NULL COMMENT '媒体类型:0-图片,1-视频',
    scene_id BIGINT UNSIGNED NOT NULL COMMENT '场景ID',
    url VARCHAR(512) NOT NULL COMMENT '媒体URL',
    metadata JSON COMMENT '元数据',
    status TINYINT UNSIGNED DEFAULT 0 COMMENT '状态:0-生成中,1-已完成,2-失败',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_scene_id (scene_id),
    INDEX idx_type (type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='媒体表';
