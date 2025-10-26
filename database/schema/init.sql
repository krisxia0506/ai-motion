-- AI-Motion PostgreSQL 数据库初始化脚本
-- 按照依赖顺序执行所有表的创建
-- 注意: 本脚本已从 MySQL 迁移至 PostgreSQL

-- 1. 小说表
CREATE TABLE IF NOT EXISTS aimotion_novel (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100),
    status SMALLINT DEFAULT 0 CHECK (status IN (0, 1, 2)),
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE aimotion_novel IS '小说表';
COMMENT ON COLUMN aimotion_novel.status IS '状态:0-草稿,1-解析中,2-已完成';
COMMENT ON COLUMN aimotion_novel.is_deleted IS '逻辑删除:0-未删除,1-已删除';

CREATE INDEX idx_status ON aimotion_novel(status);
CREATE INDEX idx_gmt_create ON aimotion_novel(gmt_create);

-- 2. 小说内容表 (大字段独立)
CREATE TABLE IF NOT EXISTS aimotion_novel_content (
    id BIGSERIAL PRIMARY KEY,
    novel_id BIGINT NOT NULL,
    content TEXT,
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_novel_id UNIQUE (novel_id)
);

COMMENT ON TABLE aimotion_novel_content IS '小说内容表';
COMMENT ON COLUMN aimotion_novel_content.novel_id IS '小说ID';
COMMENT ON COLUMN aimotion_novel_content.content IS '小说内容';

-- 3. 章节表
CREATE TABLE IF NOT EXISTS aimotion_chapter (
    id BIGSERIAL PRIMARY KEY,
    novel_id BIGINT NOT NULL,
    title VARCHAR(255),
    sequence_num INTEGER NOT NULL CHECK (sequence_num >= 0),
    content TEXT,
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_novel_sequence UNIQUE (novel_id, sequence_num)
);

COMMENT ON TABLE aimotion_chapter IS '章节表';
COMMENT ON COLUMN aimotion_chapter.is_deleted IS '逻辑删除';

CREATE INDEX idx_novel_id ON aimotion_chapter(novel_id);

-- 4. 角色表
CREATE TABLE IF NOT EXISTS aimotion_character (
    id BIGSERIAL PRIMARY KEY,
    novel_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    appearance TEXT,
    personality TEXT,
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE aimotion_character IS '角色表';
COMMENT ON COLUMN aimotion_character.is_deleted IS '逻辑删除';

CREATE INDEX idx_novel_id_char ON aimotion_character(novel_id);
CREATE INDEX idx_name ON aimotion_character(name);

-- 5. 角色图片表
CREATE TABLE IF NOT EXISTS aimotion_character_image (
    id BIGSERIAL PRIMARY KEY,
    character_id BIGINT NOT NULL,
    image_url VARCHAR(512) NOT NULL,
    image_type SMALLINT DEFAULT 0 CHECK (image_type IN (0, 1)),
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE aimotion_character_image IS '角色图片表';
COMMENT ON COLUMN aimotion_character_image.image_type IS '图片类型:0-参考图,1-场景图';
COMMENT ON COLUMN aimotion_character_image.is_deleted IS '逻辑删除';

CREATE INDEX idx_character_id ON aimotion_character_image(character_id);

-- 6. 场景表
CREATE TABLE IF NOT EXISTS aimotion_scene (
    id BIGSERIAL PRIMARY KEY,
    novel_id BIGINT NOT NULL,
    chapter_id BIGINT NOT NULL,
    sequence_num INTEGER NOT NULL CHECK (sequence_num >= 0),
    description TEXT,
    location VARCHAR(255),
    time_of_day VARCHAR(50),
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_chapter_sequence UNIQUE (chapter_id, sequence_num)
);

COMMENT ON TABLE aimotion_scene IS '场景表';
COMMENT ON COLUMN aimotion_scene.is_deleted IS '逻辑删除';

CREATE INDEX idx_novel_id_scene ON aimotion_scene(novel_id);
CREATE INDEX idx_chapter_id ON aimotion_scene(chapter_id);

-- 7. 场景角色关联表
CREATE TABLE IF NOT EXISTS aimotion_scene_character (
    id BIGSERIAL PRIMARY KEY,
    scene_id BIGINT NOT NULL,
    character_id BIGINT NOT NULL,
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uk_scene_character UNIQUE (scene_id, character_id)
);

COMMENT ON TABLE aimotion_scene_character IS '场景角色关联表';
COMMENT ON COLUMN aimotion_scene_character.is_deleted IS '逻辑删除';

CREATE INDEX idx_character_id_sc ON aimotion_scene_character(character_id);

-- 8. 媒体表
CREATE TABLE IF NOT EXISTS aimotion_media (
    id BIGSERIAL PRIMARY KEY,
    type SMALLINT NOT NULL CHECK (type IN (0, 1)),
    scene_id BIGINT NOT NULL,
    url VARCHAR(512) NOT NULL,
    metadata JSONB,
    status SMALLINT DEFAULT 0 CHECK (status IN (0, 1, 2)),
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE aimotion_media IS '媒体表';
COMMENT ON COLUMN aimotion_media.type IS '媒体类型:0-图片,1-视频';
COMMENT ON COLUMN aimotion_media.status IS '状态:0-生成中,1-已完成,2-失败';
COMMENT ON COLUMN aimotion_media.is_deleted IS '逻辑删除';

CREATE INDEX idx_scene_id ON aimotion_media(scene_id);
CREATE INDEX idx_type ON aimotion_media(type);
CREATE INDEX idx_status ON aimotion_media(status);
