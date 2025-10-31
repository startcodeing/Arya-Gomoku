# 五子棋游戏数据持久化方案

## 📋 需求分析

基于对现有代码的分析，三种对战模式需要持久化的数据包括：

### 🎮 人机对战模式 (AI Game)
- **游戏记录**：完整的游戏历史、走棋步骤
- **玩家统计**：胜率、总场次、按难度级别统计
- **AI性能数据**：搜索节点数、剪枝效率、响应时间

### 🧠 LLM对战模式 (LLM Battle)
- **游戏记录**：与LLM的对局历史、推理过程
- **LLM配置**：模型参数、API密钥、端点设置
- **对战统计**：各模型胜率、难度表现统计

### 👥 双人对战模式 (PVP)
- **房间数据**：房间信息、玩家信息、游戏状态
- **游戏记录**：完整的对战历史、走棋步骤
- **玩家档案**：昵称、胜率、历史战绩

## 🗄️ 数据库技术选型

### 推荐方案：PostgreSQL + Redis

**PostgreSQL (主数据库)**
- ✅ 支持复杂查询和关系型数据
- ✅ 优秀的JSON支持，适合存储棋盘数据
- ✅ 强一致性和ACID特性
- ✅ 丰富的数据类型（数组、JSONB等）

**Redis (缓存层)**
- ✅ 高性能键值存储
- ✅ 用于会话管理和实时数据
- ✅ 支持过期策略，适合临时数据

## 📊 数据库设计

### 1. 用户相关表

```sql
-- 用户表
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 用户会话表 (Redis实现)
-- Key: session:{session_id}
-- Value: {user_id, username, expires_at, ...}
```

### 2. 游戏记录表

```sql
-- AI游戏记录表
CREATE TABLE ai_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    difficulty VARCHAR(20) NOT NULL, -- easy, medium, hard, expert
    ai_type VARCHAR(20) NOT NULL, -- enhanced, classic
    status VARCHAR(20) NOT NULL, -- playing, human_win, ai_win, draw
    board_size INTEGER DEFAULT 15,
    move_count INTEGER DEFAULT 0,
    total_time_ms INTEGER,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- AI游戏步骤表
CREATE TABLE ai_game_moves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES ai_games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    player VARCHAR(10) NOT NULL, -- human, ai
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    timestamp_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- AI性能统计表
CREATE TABLE ai_game_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES ai_games(id) ON DELETE CASCADE,
    ai_type VARCHAR(20) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    nodes_searched BIGINT,
    cutoffs BIGINT,
    pruning_efficiency DECIMAL(5,2),
    search_time_ms INTEGER,
    move_number INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- LLM游戏记录表
CREATE TABLE llm_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    model_name VARCHAR(100) NOT NULL,
    difficulty VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    board_size INTEGER DEFAULT 15,
    move_count INTEGER DEFAULT 0,
    total_time_ms INTEGER,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- LLM游戏步骤表
CREATE TABLE llm_game_moves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES llm_games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    player VARCHAR(10) NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    reasoning TEXT,
    confidence DECIMAL(3,2),
    timestamp_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- PVP游戏记录表
CREATE TABLE pvp_games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL,
    board JSONB NOT NULL,
    move_count INTEGER DEFAULT 0,
    current_player_id UUID,
    winner_id UUID,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- PVP游戏步骤表
CREATE TABLE pvp_game_moves (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID REFERENCES pvp_games(id) ON DELETE CASCADE,
    player_id UUID NOT NULL,
    move_number INTEGER NOT NULL,
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 3. 统计和配置表

```sql
-- 用户统计表
CREATE TABLE user_statistics (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_ai_games INTEGER DEFAULT 0,
    ai_games_won INTEGER DEFAULT 0,
    ai_games_lost INTEGER DEFAULT 0,
    ai_games_drawn INTEGER DEFAULT 0,
    total_llm_games INTEGER DEFAULT 0,
    llm_games_won INTEGER DEFAULT 0,
    llm_games_lost INTEGER DEFAULT 0,
    llm_games_drawn INTEGER DEFAULT 0,
    total_pvp_games INTEGER DEFAULT 0,
    pvp_games_won INTEGER DEFAULT 0,
    pvp_games_lost INTEGER DEFAULT 0,
    pvp_games_drawn INTEGER DEFAULT 0,
    total_play_time_ms BIGINT DEFAULT 0,
    last_game_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- AI难度统计表
CREATE TABLE ai_difficulty_stats (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    difficulty VARCHAR(20) NOT NULL,
    games_played INTEGER DEFAULT 0,
    games_won INTEGER DEFAULT 0,
    games_lost INTEGER DEFAULT 0,
    games_drawn INTEGER DEFAULT 0,
    avg_response_time_ms INTEGER,
    total_play_time_ms BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, difficulty)
);

-- LLM模型统计表
CREATE TABLE llm_model_stats (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    model_name VARCHAR(100) NOT NULL,
    games_played INTEGER DEFAULT 0,
    games_won INTEGER DEFAULT 0,
    games_lost INTEGER DEFAULT 0,
    games_drawn INTEGER DEFAULT 0,
    avg_response_time_ms INTEGER,
    avg_confidence DECIMAL(3,2),
    total_play_time_ms BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, model_name)
);

-- LLM配置表
CREATE TABLE llm_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    model_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    provider VARCHAR(50) NOT NULL, -- deepseek, openai, ollama
    api_key_encrypted TEXT,
    endpoint VARCHAR(500),
    parameters JSONB,
    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, model_name)
);

-- 房间表 (活跃房间使用Redis存储)
CREATE TABLE pvp_rooms (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    creator_id UUID REFERENCES users(id),
    status VARCHAR(20) NOT NULL, -- waiting, playing, finished
    max_players INTEGER DEFAULT 2,
    current_players INTEGER DEFAULT 0,
    game_id UUID REFERENCES pvp_games(id),
    invite_code VARCHAR(20) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- PVP玩家表 (活跃玩家使用Redis存储)
CREATE TABLE pvp_room_players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id VARCHAR(100) REFERENCES pvp_rooms(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    player_name VARCHAR(100) NOT NULL,
    player_number INTEGER NOT NULL,
    is_ready BOOLEAN DEFAULT false,
    is_online BOOLEAN DEFAULT true,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_creator BOOLEAN DEFAULT false
);
```

### 4. Redis 缓存结构

```redis
# 活跃房间缓存 (TTL: 2小时)
room:{room_id} -> {room_data, players: [player_data...]}

# 用户会话缓存 (TTL: 24小时)
session:{session_id} -> {user_id, username, expires_at}

# 活跃游戏状态 (TTL: 4小时)
game_state:{game_id} -> {current_state, board, moves...}

# 聊天消息缓存 (TTL: 1小时)
chat:{room_id} -> [message_list]

# 用户在线状态 (TTL: 5分钟)
online_user:{user_id} -> {last_active, current_room}
```

## 🔧 实施计划

### Phase 1: 基础设施搭建 (1-2天)
1. **环境准备**
   - 安装PostgreSQL和Redis
   - 配置数据库连接池
   - 设置数据库迁移工具

2. **依赖管理**
   - 添加数据库驱动依赖
   - 配置环境变量

### Phase 2: 数据库结构实现 (2-3天)
1. **创建迁移脚本**
   - SQL表结构定义
   - 索引优化
   - 约束设置

2. **数据模型层**
   - 重构现有数据模型
   - 添加数据库标签
   - 实现序列化/反序列化

### Phase 3: 持久化服务层 (3-4天)
1. **仓储模式实现**
   - 游戏记录仓储
   - 用户统计仓储
   - 配置管理仓储

2. **服务层改造**
   - 游戏服务持久化
   - 统计服务实现
   - 缓存策略集成

### Phase 4: 业务逻辑集成 (2-3天)
1. **游戏流程改造**
   - 人机对战数据持久化
   - LLM对战记录存储
   - PVP数据同步

2. **统计功能实现**
   - 实时统计更新
   - 历史数据查询
   - 性能指标收集

### Phase 5: API和数据服务 (2-3天)
1. **REST API扩展**
   - 历史记录查询接口
   - 统计数据接口
   - 配置管理接口

2. **数据导入导出**
   - 游戏记录导出
   - 统计报表生成

## 🚀 技术优势

1. **高性能**：PostgreSQL + Redis双重保障
2. **可扩展**：关系型数据库支持复杂查询
3. **一致性**：ACID特性保证数据完整性
4. **灵活性**：JSONB字段适应复杂数据结构
5. **实时性**：Redis缓存支持高并发访问

## 📈 数据量预估

- **游戏记录**：1000+ 用户 × 平均10局/月 × 12月 = 120,000条/年
- **统计数据**：按用户维度存储，数据量可控
- **配置数据**：用户级配置，数据量小
- **缓存数据**：临时数据，自动过期清理

## 📁 项目结构建议

```
backend/
├── cmd/
│   └── migrate/
│       └── main.go              # 数据库迁移工具
├── internal/
│   ├── config/
│   │   ├── database.go          # 数据库配置
│   │   └── redis.go            # Redis配置
│   ├── model/
│   │   ├── persistence.go       # 持久化数据模型
│   │   └── user.go              # 用户模型
│   ├── repository/
│   │   ├── game_repository.go   # 游戏记录仓储
│   │   ├── user_repository.go   # 用户统计仓储
│   │   └── config_repository.go # 配置管理仓储
│   ├── service/
│   │   ├── persistence_service.go # 持久化服务
│   │   └── statistics_service.go  # 统计服务
│   └── database/
│       ├── migrations/           # 数据库迁移文件
│       └── seeds/              # 种子数据
├── pkg/
│   └── cache/
│       ├── redis_client.go      # Redis客户端
│       └── cache_manager.go    # 缓存管理器
└── docs/
    ├── persistence_schema.md   # 本文档
    ├── migration_guide.md      # 迁移指南
    └── api_documentation.md     # API文档
```

## 🔗 数据库依赖更新

需要在 `go.mod` 中添加以下依赖：

```go
require (
    github.com/lib/pq v1.10.9                    // PostgreSQL驱动
    github.com/redis/go-redis/v9 v9.0.5          // Redis客户端
    github.com/golang-migrate/migrate/v4 v4.15.2  // 数据库迁移工具
    gorm.io/gorm v1.24.2                          // ORM框架 (可选)
    github.com/golang-jwt/jwt/v5 v5.0.0          // JWT认证 (可选)
)
```

这个方案能够完整支持三种对战模式的数据持久化需求，同时保证良好的性能和可扩展性。