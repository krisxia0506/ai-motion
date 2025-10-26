-- PostgreSQL migration: Create media table
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

-- Comments
COMMENT ON TABLE aimotion_media IS '媒体表';
COMMENT ON COLUMN aimotion_media.id IS '主键ID';
COMMENT ON COLUMN aimotion_media.type IS '媒体类型:0-图片,1-视频';
COMMENT ON COLUMN aimotion_media.scene_id IS '场景ID';
COMMENT ON COLUMN aimotion_media.url IS '媒体URL';
COMMENT ON COLUMN aimotion_media.metadata IS '元数据';
COMMENT ON COLUMN aimotion_media.status IS '状态:0-生成中,1-已完成,2-失败';
COMMENT ON COLUMN aimotion_media.is_deleted IS '逻辑删除';
COMMENT ON COLUMN aimotion_media.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_media.gmt_modified IS '修改时间';

-- Indexes
CREATE INDEX idx_scene_id ON aimotion_media(scene_id);
CREATE INDEX idx_type ON aimotion_media(type);
CREATE INDEX idx_status ON aimotion_media(status);
