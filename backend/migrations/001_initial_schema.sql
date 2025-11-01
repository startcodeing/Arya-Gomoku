-- 001_initial_schema.sql
-- 初始数据库架构迁移

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    last_active_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户统计表
CREATE TABLE IF NOT EXISTS user_statistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
    average_game_duration_minutes DECIMAL(10,2) DEFAULT 0,
    rating INTEGER DEFAULT 1200,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建AI游戏表
CREATE TABLE IF NOT EXISTS ai_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    difficulty VARCHAR(20) NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard', 'expert')),
    board_size INTEGER DEFAULT 15 CHECK (board_size BETWEEN 9 AND 19),
    moves JSONB NOT NULL DEFAULT '[]',
    result VARCHAR(20) CHECK (result IN ('user_win', 'ai_win', 'draw', 'ongoing')),
    winner VARCHAR(20) CHECK (winner IN ('user', 'ai', 'draw')),
    duration_seconds INTEGER,
    move_count INTEGER DEFAULT 0,
    user_color VARCHAR(10) DEFAULT 'black' CHECK (user_color IN ('black', 'white')),
    ai_thinking_time_ms INTEGER DEFAULT 0,
    game_metadata JSONB DEFAULT '{}',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建LLM游戏表
CREATE TABLE IF NOT EXISTS llm_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    model_name VARCHAR(100) NOT NULL,
    board_size INTEGER DEFAULT 15 CHECK (board_size BETWEEN 9 AND 19),
    moves JSONB NOT NULL DEFAULT '[]',
    result VARCHAR(20) CHECK (result IN ('user_win', 'llm_win', 'draw', 'ongoing')),
    winner VARCHAR(20) CHECK (winner IN ('user', 'llm', 'draw')),
    duration_seconds INTEGER,
    move_count INTEGER DEFAULT 0,
    user_color VARCHAR(10) DEFAULT 'black' CHECK (user_color IN ('black', 'white')),
    llm_response_time_ms INTEGER DEFAULT 0,
    game_metadata JSONB DEFAULT '{}',
    conversation_history JSONB DEFAULT '[]',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建PVP游戏表
CREATE TABLE IF NOT EXISTS pvp_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id VARCHAR(100) NOT NULL,
    player1_id UUID REFERENCES users(id) ON DELETE SET NULL,
    player2_id UUID REFERENCES users(id) ON DELETE SET NULL,
    board_size INTEGER DEFAULT 15 CHECK (board_size BETWEEN 9 AND 19),
    moves JSONB NOT NULL DEFAULT '[]',
    result VARCHAR(20) CHECK (result IN ('player1_win', 'player2_win', 'draw', 'ongoing', 'abandoned')),
    winner_id UUID REFERENCES users(id) ON DELETE SET NULL,
    duration_seconds INTEGER,
    move_count INTEGER DEFAULT 0,
    current_player VARCHAR(10) DEFAULT 'black' CHECK (current_player IN ('black', 'white')),
    game_metadata JSONB DEFAULT '{}',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建PVP房间玩家表
CREATE TABLE IF NOT EXISTS pvp_room_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES pvp_games(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    color VARCHAR(10) NOT NULL CHECK (color IN ('black', 'white')),
    is_ready BOOLEAN DEFAULT false,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(game_id, user_id),
    UNIQUE(game_id, color)
);

-- 创建LLM配置表
CREATE TABLE IF NOT EXISTS llm_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_name VARCHAR(100) UNIQUE NOT NULL,
    api_endpoint VARCHAR(500) NOT NULL,
    api_key_encrypted VARCHAR(500),
    max_tokens INTEGER DEFAULT 1000,
    temperature DECIMAL(3,2) DEFAULT 0.7 CHECK (temperature BETWEEN 0 AND 2),
    system_prompt TEXT,
    is_active BOOLEAN DEFAULT true,
    rate_limit_per_minute INTEGER DEFAULT 60,
    timeout_seconds INTEGER DEFAULT 30,
    config_metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
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

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为所有表创建更新时间触发器
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_sessions_updated_at BEFORE UPDATE ON user_sessions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_statistics_updated_at BEFORE UPDATE ON user_statistics FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ai_games_updated_at BEFORE UPDATE ON ai_games FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_llm_games_updated_at BEFORE UPDATE ON llm_games FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_pvp_games_updated_at BEFORE UPDATE ON pvp_games FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_pvp_room_players_updated_at BEFORE UPDATE ON pvp_room_players FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_llm_configs_updated_at BEFORE UPDATE ON llm_configs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();