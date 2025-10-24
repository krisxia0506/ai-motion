-- AI-Motion 数据库初始化脚本
-- 按照依赖顺序执行所有表的创建

-- 1. 小说表
CREATE TABLE IF NOT EXISTS aimotion_novel (
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

-- 2. 小说内容表 (大字段独立)
CREATE TABLE IF NOT EXISTS aimotion_novel_content (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    novel_id BIGINT UNSIGNED NOT NULL COMMENT '小说ID',
    content LONGTEXT COMMENT '小说内容',
    
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    UNIQUE KEY uk_novel_id (novel_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小说内容表';

-- 3. 章节表
CREATE TABLE IF NOT EXISTS aimotion_chapter (
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

-- 4. 角色表
CREATE TABLE IF NOT EXISTS aimotion_character (
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

-- 5. 角色图片表
CREATE TABLE IF NOT EXISTS aimotion_character_image (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    character_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    image_url VARCHAR(512) NOT NULL COMMENT '图片URL',
    image_type TINYINT UNSIGNED DEFAULT 0 COMMENT '图片类型:0-参考图,1-场景图',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    INDEX idx_character_id (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色图片表';

-- 6. 场景表
CREATE TABLE IF NOT EXISTS aimotion_scene (
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

-- 7. 场景角色关联表
CREATE TABLE IF NOT EXISTS aimotion_scene_character (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    scene_id BIGINT UNSIGNED NOT NULL COMMENT '场景ID',
    character_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    
    is_deleted TINYINT UNSIGNED DEFAULT 0 COMMENT '逻辑删除',
    gmt_create DATETIME NOT NULL COMMENT '创建时间',
    gmt_modified DATETIME NOT NULL COMMENT '修改时间',
    
    UNIQUE KEY uk_scene_character (scene_id, character_id),
    INDEX idx_character_id (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='场景角色关联表';

-- 8. 媒体表
CREATE TABLE IF NOT EXISTS aimotion_media (
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
