-- PostgreSQL migration: Create tasks table
CREATE TABLE IF NOT EXISTS aimotion_task (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    novel_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    progress_step VARCHAR(100),
    progress_step_index INT DEFAULT 0,
    progress_percentage INT DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    progress_details JSONB,
    error_code INT,
    error_message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    failed_at TIMESTAMP NULL
);

-- Comments
COMMENT ON TABLE aimotion_task IS '异步任务表';
COMMENT ON COLUMN aimotion_task.id IS '主键ID';
COMMENT ON COLUMN aimotion_task.user_id IS '用户ID，关联 Supabase Auth 用户';
COMMENT ON COLUMN aimotion_task.novel_id IS '关联的小说ID';
COMMENT ON COLUMN aimotion_task.status IS '任务状态:pending-待处理,processing-处理中,completed-已完成,failed-失败,cancelled-已取消';
COMMENT ON COLUMN aimotion_task.progress_step IS '当前步骤描述';
COMMENT ON COLUMN aimotion_task.progress_step_index IS '步骤索引（1-6）';
COMMENT ON COLUMN aimotion_task.progress_percentage IS '进度百分比（0-100）';
COMMENT ON COLUMN aimotion_task.progress_details IS '进度详情（字符、场景生成数量）';
COMMENT ON COLUMN aimotion_task.error_code IS '错误代码';
COMMENT ON COLUMN aimotion_task.error_message IS '错误信息';
COMMENT ON COLUMN aimotion_task.created_at IS '创建时间';
COMMENT ON COLUMN aimotion_task.updated_at IS '修改时间';
COMMENT ON COLUMN aimotion_task.completed_at IS '完成时间';
COMMENT ON COLUMN aimotion_task.failed_at IS '失败时间';

-- Indexes
CREATE INDEX idx_task_user_id ON aimotion_task(user_id);
CREATE INDEX idx_task_novel_id ON aimotion_task(novel_id);
CREATE INDEX idx_task_status ON aimotion_task(status);
CREATE INDEX idx_task_created_at ON aimotion_task(created_at DESC);
CREATE INDEX idx_task_user_status ON aimotion_task(user_id, status);

-- Trigger for automatic updated_at
CREATE OR REPLACE FUNCTION update_aimotion_task_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_aimotion_task_updated_at
BEFORE UPDATE ON aimotion_task
FOR EACH ROW
EXECUTE FUNCTION update_aimotion_task_updated_at();
