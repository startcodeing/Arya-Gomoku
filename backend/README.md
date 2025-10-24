# 五子棋后端服务 (Gomoku Backend)

基于 Golang + Gin 框架的五子棋游戏后端API服务。

## 环境要求

- **Golang**: 1.21 或更高版本
- **操作系统**: Windows/macOS/Linux

## 安装 Golang

### Windows
1. 访问 [Golang官网](https://golang.org/dl/)
2. 下载 Windows 安装包 (go1.21.x.windows-amd64.msi)
3. 运行安装包，按照向导完成安装
4. 验证安装：
   ```cmd
   go version
   ```

### macOS
```bash
# 使用 Homebrew
brew install go

# 或者下载官方安装包
# https://golang.org/dl/
```

### Linux (Ubuntu/Debian)
```bash
# 方法1: 使用包管理器
sudo apt update
sudo apt install golang-go

# 方法2: 下载官方二进制包
wget https://golang.org/dl/go1.21.x.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.x.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

## 项目设置

### 1. 安装依赖
```bash
cd backend
go mod tidy
```

### 2. 启动服务
```bash
go run main.go
```

服务将在 http://localhost:8080 启动

### 3. 构建可执行文件
```bash
# Windows
go build -o gomoku-server.exe main.go

# macOS/Linux
go build -o gomoku-server main.go
```

## API 接口文档

### 基础信息
- **基础URL**: http://localhost:8080/api
- **Content-Type**: application/json
- **CORS**: 支持跨域请求 (允许 http://localhost:5173)

### AI 对战接口

#### 1. 获取AI移动
```http
POST /api/ai/move
Content-Type: application/json

{
  "board": [
    [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
    [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],
    ...
  ],
  "currentPlayer": 2
}
```

**响应**:
```json
{
  "x": 7,
  "y": 7,
  "player": 2,
  "gameStatus": "playing",
  "winner": 0
}
```

#### 2. 获取游戏状态
```http
GET /api/ai/status
```

**响应**:
```json
{
  "message": "AI service is running",
  "status": "active"
}
```

#### 3. 重置游戏
```http
POST /api/ai/reset
```

**响应**:
```json
{
  "message": "Game reset successfully"
}
```

### 在线匹配接口 (预留功能)

#### 1. 开始匹配
```http
POST /api/match/start
Content-Type: application/json

{
  "playerId": "player123",
  "gameMode": "pvp"
}
```

#### 2. 加入房间
```http
POST /api/match/join
Content-Type: application/json

{
  "roomId": "room123",
  "playerId": "player456"
}
```

#### 3. 获取房间状态
```http
GET /api/match/status/{roomId}
```

#### 4. 在房间中移动
```http
POST /api/match/{roomId}/move
Content-Type: application/json

{
  "playerId": "player123",
  "x": 7,
  "y": 8
}
```

#### 5. 离开房间
```http
POST /api/match/{roomId}/leave
Content-Type: application/json

{
  "playerId": "player123"
}
```

#### 6. 获取活跃房间数
```http
GET /api/match/rooms
```

## 数据模型

### 棋盘表示
- **大小**: 15x15
- **值含义**: 
  - `0`: 空位
  - `1`: 玩家棋子
  - `2`: AI棋子

### 游戏状态
- `playing`: 游戏进行中
- `human_win`: 玩家获胜
- `ai_win`: AI获胜
- `draw`: 平局

## AI 算法

### 启发式评估函数
```go
// 评分标准
const (
    WIN_SCORE     = 10000  // 连五
    FOUR_SCORE    = 1000   // 活四
    BLOCK_FOUR    = 100    // 冲四
    THREE_SCORE   = 100    // 活三
    BLOCK_THREE   = 10     // 眠三
    TWO_SCORE     = 10     // 活二
    BLOCK_TWO     = 1      // 眠二
)
```

### 决策优先级
1. **获胜移动**: 如果能连成五子，立即获胜
2. **阻止对手获胜**: 阻止对手连成五子
3. **创造威胁**: 形成活四或活三
4. **防守威胁**: 阻止对手形成威胁
5. **战略位置**: 在重要位置落子
6. **随机选择**: 在有效位置中随机选择

## 项目结构

```
backend/
├── main.go                    # 主程序入口
├── go.mod                     # Go模块文件
├── go.sum                     # 依赖校验文件
└── internal/                  # 内部包
    ├── model/                 # 数据模型
    │   └── board.go          # 棋盘和游戏逻辑
    ├── service/              # 业务服务
    │   ├── ai_service.go     # AI算法实现
    │   └── game_service.go   # 游戏管理服务
    └── controller/           # HTTP控制器
        ├── ai_controller.go  # AI接口控制器
        └── game_controller.go # 游戏房间控制器
```

## 开发指南

### 添加新的AI策略
1. 在 `ai_service.go` 中修改 `evaluatePosition` 函数
2. 调整评分权重和优先级
3. 测试新策略的效果

### 添加新的API接口
1. 在相应的 controller 中添加处理函数
2. 在 `main.go` 中注册路由
3. 更新API文档

### 错误处理
- 使用标准的HTTP状态码
- 返回结构化的错误信息
- 记录详细的错误日志

## 测试

### 运行测试
```bash
go test ./...
```

### 手动测试API
```bash
# 测试AI移动接口
curl -X POST http://localhost:8080/api/ai/move \
  -H "Content-Type: application/json" \
  -d '{
    "board": [[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,1,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]],
    "currentPlayer": 2
  }'
```

## 部署

### 生产环境部署
```bash
# 构建
go build -ldflags="-s -w" -o gomoku-server main.go

# 运行
./gomoku-server
```

### Docker 部署 (可选)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o gomoku-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gomoku-server .
EXPOSE 8080
CMD ["./gomoku-server"]
```

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # Windows
   netstat -ano | findstr :8080
   
   # macOS/Linux
   lsof -i :8080
   ```

2. **CORS 错误**
   - 确保前端运行在 http://localhost:5173
   - 检查 CORS 配置是否正确

3. **依赖问题**
   ```bash
   go clean -modcache
   go mod tidy
   ```

## 性能优化

- 使用连接池管理数据库连接
- 实现请求缓存机制
- 优化AI算法的搜索深度
- 添加请求限流和防护

---

**后端服务已准备就绪！** 🚀