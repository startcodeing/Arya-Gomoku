-- 002_initial_data_sqlite.sql
-- SQLite兼容的初始数据迁移

-- 插入默认管理员用户
INSERT OR REPLACE INTO users (
    id, username, email, password_hash, nickname, role, is_active, email_verified, created_at, updated_at
) VALUES (
    'admin-001', 
    'admin', 
    'admin@gomoku.local', 
    '$2a$10$nwkc.2ldZ8yzYUxnSATQEO11YRa9rK0rgYpDZ5qSJWTNcWjHsFKtC', -- password123
    'Administrator',
    'admin',
    1,
    1,
    datetime('now'),
    datetime('now')
);

-- 插入测试用户
INSERT OR REPLACE INTO users (
    id, username, email, password_hash, nickname, role, is_active, email_verified, created_at, updated_at
) VALUES (
    'user-001', 
    'testuser', 
    'test@gomoku.local', 
    '$2a$10$nwkc.2ldZ8yzYUxnSATQEO11YRa9rK0rgYpDZ5qSJWTNcWjHsFKtC', -- password123
    'Test User',
    'user',
    1,
    1,
    datetime('now'),
    datetime('now')
);

-- 插入用户统计数据
INSERT OR IGNORE INTO user_statistics (
    id, user_id, total_games, total_wins, total_losses, total_draws,
    ai_games, ai_wins, llm_games, llm_wins, pvp_games, pvp_wins,
    win_streak, max_win_streak, total_playtime_minutes, average_game_duration_minutes,
    rating, created_at, updated_at
) VALUES 
(
    'stats-admin-001', 'admin-001', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1200,
    datetime('now'), datetime('now')
),
(
    'stats-user-001', 'user-001', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1200,
    datetime('now'), datetime('now')
);

-- 插入LLM配置
INSERT OR IGNORE INTO llm_configs (
    id, model_name, api_endpoint, max_tokens, temperature, system_prompt, 
    is_active, rate_limit_per_minute, timeout_seconds, config_metadata, created_at, updated_at
) VALUES 
(
    'llm-config-001',
    'gpt-3.5-turbo',
    'https://api.openai.com/v1/chat/completions',
    1000,
    0.7,
    'You are an expert Gomoku (Five in a Row) player. You will play against a human opponent on a 15x15 board. Your goal is to get five stones in a row (horizontally, vertically, or diagonally) before your opponent does. Analyze the board position carefully and make strategic moves. Respond with only the coordinates in the format "x,y" where x and y are numbers from 0 to 14.',
    1,
    60,
    30,
    '{"provider": "openai", "version": "v1"}',
    datetime('now'),
    datetime('now')
),
(
    'llm-config-002',
    'claude-3-sonnet',
    'https://api.anthropic.com/v1/messages',
    1000,
    0.7,
    'You are an expert Gomoku (Five in a Row) player. You will play against a human opponent on a 15x15 board. Your goal is to get five stones in a row (horizontally, vertically, or diagonally) before your opponent does. Analyze the board position carefully and make strategic moves. Respond with only the coordinates in the format "x,y" where x and y are numbers from 0 to 14.',
    1,
    60,
    30,
    '{"provider": "anthropic", "version": "v1"}',
    datetime('now'),
    datetime('now')
);

-- 插入示例AI游戏记录
INSERT OR IGNORE INTO ai_games (
    id, user_id, difficulty, board_size, moves, result, winner, 
    duration_seconds, move_count, user_color, ai_thinking_time_ms, 
    game_metadata, started_at, finished_at, created_at, updated_at
) VALUES (
    'ai-game-001',
    'user-001',
    'medium',
    15,
    '[{"x":7,"y":7,"player":"black"},{"x":8,"y":8,"player":"white"},{"x":6,"y":6,"player":"black"},{"x":9,"y":9,"player":"white"},{"x":5,"y":5,"player":"black"}]',
    'user_win',
    'user',
    300,
    5,
    'black',
    2500,
    '{"opening": "center", "strategy": "offensive"}',
    datetime('now', '-1 day'),
    datetime('now', '-1 day', '+5 minutes'),
    datetime('now', '-1 day'),
    datetime('now', '-1 day')
);

-- 插入示例LLM游戏记录
INSERT OR IGNORE INTO llm_games (
    id, user_id, model_name, board_size, moves, result, winner,
    duration_seconds, move_count, user_color, llm_response_time_ms,
    game_metadata, conversation_history, started_at, finished_at, created_at, updated_at
) VALUES (
    'llm-game-001',
    'user-001',
    'gpt-3.5-turbo',
    15,
    '[{"x":7,"y":7,"player":"black"},{"x":8,"y":7,"player":"white"},{"x":6,"y":7,"player":"black"},{"x":9,"y":7,"player":"white"},{"x":5,"y":7,"player":"black"}]',
    'user_win',
    'user',
    420,
    5,
    'black',
    3200,
    '{"opening": "center", "strategy": "defensive"}',
    '[{"role":"user","content":"I played at 7,7"},{"role":"assistant","content":"8,7"},{"role":"user","content":"I played at 6,7"},{"role":"assistant","content":"9,7"},{"role":"user","content":"I played at 5,7"}]',
    datetime('now', '-2 hours'),
    datetime('now', '-2 hours', '+7 minutes'),
    datetime('now', '-2 hours'),
    datetime('now', '-2 hours')
);