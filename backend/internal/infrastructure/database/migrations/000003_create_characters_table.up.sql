-- PostgreSQL migration: Create characters table
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

-- Comments
COMMENT ON TABLE aimotion_character IS '角色表';
COMMENT ON COLUMN aimotion_character.id IS '主键ID';
COMMENT ON COLUMN aimotion_character.novel_id IS '小说ID';
COMMENT ON COLUMN aimotion_character.name IS '角色名称';
COMMENT ON COLUMN aimotion_character.appearance IS '外貌描述';
COMMENT ON COLUMN aimotion_character.personality IS '性格描述';
COMMENT ON COLUMN aimotion_character.is_deleted IS '逻辑删除';
COMMENT ON COLUMN aimotion_character.gmt_create IS '创建时间';
COMMENT ON COLUMN aimotion_character.gmt_modified IS '修改时间';

-- Indexes
CREATE INDEX idx_novel_id ON aimotion_character(novel_id);
CREATE INDEX idx_name ON aimotion_character(name);
