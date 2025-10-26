-- 创建 tasks 表
CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID，关联 Supabase Auth 用户',
    novel_id VARCHAR(36) COMMENT '关联的小说ID',
    status ENUM('pending', 'processing', 'completed', 'failed', 'cancelled') NOT NULL DEFAULT 'pending',
    progress_step VARCHAR(100) COMMENT '当前步骤描述',
    progress_step_index INT DEFAULT 0 COMMENT '步骤索引（1-6）',
    progress_percentage INT DEFAULT 0 COMMENT '进度百分比（0-100）',
    progress_details JSON COMMENT '进度详情（字符、场景生成数量）',
    error_code INT COMMENT '错误代码',
    error_message TEXT COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    failed_at TIMESTAMP NULL COMMENT '失败时间',

    INDEX idx_user_id (user_id),
    INDEX idx_novel_id (novel_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at DESC),
    INDEX idx_user_status (user_id, status) COMMENT '组合索引，用于按用户和状态查询'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='异步任务表';
