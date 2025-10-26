-- PostgreSQL migration: Create chapters table
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

-- Comments
COMMENT ON TABLE aimotion_chapter IS '章节表';
COMMENT ON COLUMN aimotion_chapter.id IS '主键ID';
COMMENT ON COLUMN aimotion_chapter.novel_id IS '小说ID';
COMMENT ON COLUMN aimotion_chapter.title IS '章节标题';
COMMENT ON COLUMN aimotion_chapter.sequence_num IS '章节序号';
COMMENT ON COLUMN aimotion_chapter.content IS '章节内容';
COMMENT ON COLUMN aimotion_chapter.is_deleted IS '逻辑删除';
COMMENT ON COLUMN aimotion_chapter.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_chapter.gmt_modified IS '修改时间';

-- Indexes
CREATE INDEX idx_novel_id ON aimotion_chapter(novel_id);
