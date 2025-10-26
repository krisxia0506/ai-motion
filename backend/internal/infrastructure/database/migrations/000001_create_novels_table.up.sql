-- PostgreSQL migration: Create novels table
CREATE TABLE IF NOT EXISTS aimotion_novel (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100),
    status SMALLINT DEFAULT 0 CHECK (status IN (0, 1, 2)),
    is_deleted SMALLINT DEFAULT 0 CHECK (is_deleted IN (0, 1)),
    gmt_create TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    gmt_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Comments
COMMENT ON TABLE aimotion_novel IS '小说表';
COMMENT ON COLUMN aimotion_novel.id IS '主键ID';
COMMENT ON COLUMN aimotion_novel.title IS '小说标题';
COMMENT ON COLUMN aimotion_novel.author IS '作者';
COMMENT ON COLUMN aimotion_novel.status IS '状态:0-草稿,1-解析中,2-已完成';
COMMENT ON COLUMN aimotion_novel.is_deleted IS '逻辑删除:0-未删除,1-已删除';
COMMENT ON COLUMN aimotion_novel.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_novel.gmt_modified IS '修改时间';

-- Indexes
CREATE INDEX idx_status ON aimotion_novel(status);
CREATE INDEX idx_gmt_create ON aimotion_novel(gmt_create);
