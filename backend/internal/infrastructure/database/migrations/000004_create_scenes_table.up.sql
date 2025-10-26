-- PostgreSQL migration: Create scenes table
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

-- Comments
COMMENT ON TABLE aimotion_scene IS '场景表';
COMMENT ON COLUMN aimotion_scene.id IS '主键ID';
COMMENT ON COLUMN aimotion_scene.novel_id IS '小说ID';
COMMENT ON COLUMN aimotion_scene.chapter_id IS '章节ID';
COMMENT ON COLUMN aimotion_scene.sequence_num IS '场景序号';
COMMENT ON COLUMN aimotion_scene.description IS '场景描述';
COMMENT ON COLUMN aimotion_scene.location IS '地点';
COMMENT ON COLUMN aimotion_scene.time_of_day IS '时间段';
COMMENT ON COLUMN aimotion_scene.is_deleted IS '逻辑删除';
COMMENT ON COLUMN aimotion_scene.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_scene.gmt_modified IS '修改时间';

-- Indexes
CREATE INDEX idx_novel_id ON aimotion_scene(novel_id);
CREATE INDEX idx_chapter_id ON aimotion_scene(chapter_id);
