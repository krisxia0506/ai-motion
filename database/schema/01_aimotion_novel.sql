-- 小说表
CREATE TABLE aimotion_novel (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    title VARCHAR(255) NOT NULL COMMENT '小说标题',
    author VARCHAR(100) COMMENT '作者',
    status TINYINT UNSIGNED DEFAULT 0 COMMENT '状态:0-草稿,1-解析中,2-已完成',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除:0-未删除,1-已删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_status (status),
    INDEX idx_gmt_create (gmt_create)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小说表';
