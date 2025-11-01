-- 001_initial_schema_sqlite.sql
-- SQLite兼容的初始数据库架构迁移

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    nickname TEXT,
    avatar_url TEXT,
    role TEXT DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
    is_active INTEGER DEFAULT 1,
    email_verified INTEGER DEFAULT 0,
    last_active_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token TEXT UNIQUE NOT NULL,
    refresh_token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户统计表
CREATE TABLE IF NOT EXISTS user_statistics (
    id TEXT PRIMARY KEY,
    user_id TEXT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_games INTEGER DEFAULT 0,
    total_wins INTEGER DEFAULT 0,
    total_losses INTEGER DEFAULT 0,
    total_draws INTEGER DEFAULT 0,
    ai_games INTEGER DEFAULT 0,
    ai_wins INTEGER DEFAULT 0,
    llm_games INTEGER DEFAULT 0,
    llm_wins INTEGER DEFAULT 0,
    pvp_games INTEGER DEFAULT 0,
    pvp_wins INTEGER DEFAULT 0,
    win_streak INTEGER DEFAULT 0,
    max_win_streak INTEGER DEFAULT 0,
    total_playtime_minutes INTEGER DEFAULT 0,
    average_game_duration_minutes REAL DEFAULT 0,
    rating INTEGER DEFAULT 1200,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建AI游戏表
CREATE TABLE IF NOT EXISTS ai_games (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    difficulty TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard', 'expert')),
    board_size INTEGER DEFAULT 15 CHECK (board_size BETWEEN 9 AND 19),
    moves TEXT NOT NULL DEFAULT '[]',
    result TEXT CHECK (result IN ('user_win', 'ai_win', 'draw', 'ongoing')),
    winner TEXT CHECK (winner IN ('user', 'ai', 'draw')),
    duration_seconds INTEGER,
    move_count INTEGER DEFAULT 0,
    user_color TEXT DEFAULT 'black' CHECK (user_color IN ('black', 'white')),
    ai_thinking_time_ms INTEGER DEFAULT 0,
    game_metadata TEXT DEFAULT '{}',
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建LLM游戏表
CREATE TABLE IF NOT EXISTS llm_games (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    model_name TEXT NOT NULL,
    board_size INTEGER DEFAULT 15 CHECK (board_size BETWEEN 9 AND 19),
    moves TEXT NOT NULL DEFAULT '[]',
    result TEXT CHECK (result IN ('user_win', 'llm_win', 'draw', 'ongoing')),
    winner TEXT CHECK (winner IN ('user', 'llm', 'draw')),
    duration_seconds INTEGER,
    move_count INTEGER DEFAULT 0,
    user_color TEXT DEFAULT 'black' CHECK (user_color IN ('black', 'white')),
    llm_response_time_ms INTEGER DEFAULT 0,
    game_metadata TEXT DEFAULT '{}',
    conversation_history TEXT DEFAULT '[]',
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建PVP游戏表
CREATE TABLE IF NOT EXISTS pvp_games (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    player1_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    player2_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    board_size INTEGER DEFAULT 15 CHECK (board_size BETWEEN 9 AND 19),
    moves TEXT NOT NULL DEFAULT '[]',
    result TEXT CHECK (result IN ('player1_win', 'player2_win', 'draw', 'ongoing', 'abandoned')),
    winner_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    duration_seconds INTEGER,
    move_count INTEGER DEFAULT 0,
    current_player TEXT DEFAULT 'black' CHECK (current_player IN ('black', 'white')),
    game_metadata TEXT DEFAULT '{}',
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    finished_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建PVP房间玩家表
CREATE TABLE IF NOT EXISTS pvp_room_players (
    id TEXT PRIMARY KEY,
    game_id TEXT NOT NULL REFERENCES pvp_games(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    color TEXT NOT NULL CHECK (color IN ('black', 'white')),
    is_ready INTEGER DEFAULT 0,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    left_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, user_id),
    UNIQUE(game_id, color)
);

-- 创建LLM配置表
CREATE TABLE IF NOT EXISTS llm_configs (
    id TEXT PRIMARY KEY,
    model_name TEXT UNIQUE NOT NULL,
    api_endpoint TEXT NOT NULL,
    api_key_encrypted TEXT,
    max_tokens INTEGER DEFAULT 1000,
    temperature REAL DEFAULT 0.7 CHECK (temperature BETWEEN 0 AND 2),
    system_prompt TEXT,
    is_active INTEGER DEFAULT 1,
    rate_limit_per_minute INTEGER DEFAULT 60,
    timeout_seconds INTEGER DEFAULT 30,
    config_metadata TEXT DEFAULT '{}',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_last_active_at ON users(last_active_at);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_session_token ON user_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_user_statistics_user_id ON user_statistics(user_id);
CREATE INDEX IF NOT EXISTS idx_user_statistics_rating ON user_statistics(rating);
CREATE INDEX IF NOT EXISTS idx_user_statistics_total_games ON user_statistics(total_games);

CREATE INDEX IF NOT EXISTS idx_ai_games_user_id ON ai_games(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_games_difficulty ON ai_games(difficulty);
CREATE INDEX IF NOT EXISTS idx_ai_games_result ON ai_games(result);
CREATE INDEX IF NOT EXISTS idx_ai_games_started_at ON ai_games(started_at);
CREATE INDEX IF NOT EXISTS idx_ai_games_finished_at ON ai_games(finished_at);

CREATE INDEX IF NOT EXISTS idx_llm_games_user_id ON llm_games(user_id);
CREATE INDEX IF NOT EXISTS idx_llm_games_model_name ON llm_games(model_name);
CREATE INDEX IF NOT EXISTS idx_llm_games_result ON llm_games(result);
CREATE INDEX IF NOT EXISTS idx_llm_games_started_at ON llm_games(started_at);
CREATE INDEX IF NOT EXISTS idx_llm_games_finished_at ON llm_games(finished_at);

CREATE INDEX IF NOT EXISTS idx_pvp_games_room_id ON pvp_games(room_id);
CREATE INDEX IF NOT EXISTS idx_pvp_games_player1_id ON pvp_games(player1_id);
CREATE INDEX IF NOT EXISTS idx_pvp_games_player2_id ON pvp_games(player2_id);
CREATE INDEX IF NOT EXISTS idx_pvp_games_winner_id ON pvp_games(winner_id);
CREATE INDEX IF NOT EXISTS idx_pvp_games_result ON pvp_games(result);
CREATE INDEX IF NOT EXISTS idx_pvp_games_started_at ON pvp_games(started_at);
CREATE INDEX IF NOT EXISTS idx_pvp_games_finished_at ON pvp_games(finished_at);

CREATE INDEX IF NOT EXISTS idx_pvp_room_players_game_id ON pvp_room_players(game_id);
CREATE INDEX IF NOT EXISTS idx_pvp_room_players_user_id ON pvp_room_players(user_id);
CREATE INDEX IF NOT EXISTS idx_pvp_room_players_joined_at ON pvp_room_players(joined_at);

CREATE INDEX IF NOT EXISTS idx_llm_configs_model_name ON llm_configs(model_name);
CREATE INDEX IF NOT EXISTS idx_llm_configs_is_active ON llm_configs(is_active);