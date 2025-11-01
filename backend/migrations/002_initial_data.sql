-- 002_initial_data.sql
-- 初始数据插入

-- 插入默认管理员用户（如果不存在）
INSERT INTO users (id, username, email, password_hash, nickname, role, is_active, email_verified)
SELECT 
    gen_random_uuid(),
    'admin',
    'admin@arya-gomoku.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
    '系统管理员',
    'admin',
    true,
    true
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE username = 'admin' OR email = 'admin@arya-gomoku.com'
);

-- 为管理员用户创建统计记录
INSERT INTO user_statistics (user_id)
SELECT id FROM users WHERE username = 'admin'
ON CONFLICT (user_id) DO NOTHING;

-- 插入默认LLM配置
INSERT INTO llm_configs (model_name, api_endpoint, max_tokens, temperature, system_prompt, is_active, rate_limit_per_minute, timeout_seconds)
VALUES 
    (
        'gpt-3.5-turbo',
        'https://api.openai.com/v1/chat/completions',
        1000,
        0.7,
        '你是一个五子棋AI助手。你需要分析当前棋盘状态，并选择最佳的下棋位置。请以坐标形式回复，例如：(8,8)。',
        true,
        60,
        30
    ),
    (
        'gpt-4',
        'https://api.openai.com/v1/chat/completions',
        1000,
        0.5,
        '你是一个专业的五子棋AI。请分析棋盘局势，考虑攻防平衡，选择最优落子位置。回复格式：(x,y)。',
        true,
        30,
        45
    ),
    (
        'claude-3-sonnet',
        'https://api.anthropic.com/v1/messages',
        1000,
        0.6,
        '你是五子棋专家。分析当前棋局，选择最佳落子点。考虑：1)形成五子连线 2)阻止对手连线 3)创造多重威胁。回复坐标：(x,y)。',
        true,
        40,
        40
    )
ON CONFLICT (model_name) DO NOTHING;

-- 创建一些示例游戏数据（可选，用于测试）
-- 注意：这些数据仅用于开发和测试环境

-- 插入测试用户（如果在开发环境）
INSERT INTO users (username, email, password_hash, nickname, is_active, email_verified)
SELECT 
    'testuser',
    'test@example.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: password
    '测试用户',
    true,
    true
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE username = 'testuser' OR email = 'test@example.com'
);

-- 为测试用户创建统计记录
INSERT INTO user_statistics (user_id, total_games, total_wins, ai_games, ai_wins, rating)
SELECT 
    id,
    10,
    6,
    8,
    5,
    1350
FROM users WHERE username = 'testuser'
ON CONFLICT (user_id) DO NOTHING;

-- 插入一些示例AI游戏记录
INSERT INTO ai_games (user_id, difficulty, moves, result, winner, duration_seconds, move_count, finished_at)
SELECT 
    u.id,
    'medium',
    '[{"x":7,"y":7,"player":"user"},{"x":8,"y":8,"player":"ai"},{"x":6,"y":6,"player":"user"},{"x":9,"y":9,"player":"ai"},{"x":5,"y":5,"player":"user"}]'::jsonb,
    'user_win',
    'user',
    180,
    9,
    CURRENT_TIMESTAMP - INTERVAL '1 day'
FROM users u WHERE u.username = 'testuser'
UNION ALL
SELECT 
    u.id,
    'hard',
    '[{"x":7,"y":7,"player":"user"},{"x":8,"y":7,"player":"ai"},{"x":7,"y":8,"player":"user"},{"x":8,"y":8,"player":"ai"}]'::jsonb,
    'ai_win',
    'ai',
    240,
    12,
    CURRENT_TIMESTAMP - INTERVAL '2 days'
FROM users u WHERE u.username = 'testuser';

-- 插入一些示例LLM游戏记录
INSERT INTO llm_games (user_id, model_name, moves, result, winner, duration_seconds, move_count, finished_at)
SELECT 
    u.id,
    'gpt-3.5-turbo',
    '[{"x":7,"y":7,"player":"user"},{"x":8,"y":8,"player":"llm"},{"x":6,"y":6,"player":"user"}]'::jsonb,
    'ongoing',
    NULL,
    NULL,
    3,
    NULL
FROM users u WHERE u.username = 'testuser';

-- 创建数据库版本表（用于跟踪迁移）
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 记录已应用的迁移
INSERT INTO schema_migrations (version) VALUES ('001_initial_schema') ON CONFLICT DO NOTHING;
INSERT INTO schema_migrations (version) VALUES ('002_initial_data') ON CONFLICT DO NOTHING;