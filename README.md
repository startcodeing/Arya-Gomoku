# 五子棋人机对战系统 (Gomoku AI vs Human)

一个基于 Golang + Vue3 的五子棋人机对战游戏，支持智能AI算法和在线匹配功能扩展。

## 项目特性

- 🎮 **人机对战**: 与智能AI进行五子棋对战
- 🧠 **启发式AI**: 采用优先级算法（赢棋 > 防守 > 进攻 > 随机）
- 🎨 **现代界面**: 基于Vue3的响应式游戏界面
- 🔄 **实时交互**: 流畅的游戏体验和状态反馈
- 📊 **游戏统计**: 记录胜负数据和游戏历史
- 🌐 **扩展架构**: 预留在线匹配功能的完整架构

## 技术栈

### 后端 (Backend)
- **语言**: Golang 1.21+
- **框架**: Gin Web Framework
- **架构**: RESTful API
- **CORS**: 支持跨域请求

### 前端 (Frontend)
- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **HTTP客户端**: Axios
- **样式**: CSS3 + 响应式设计

## 项目结构

```
Gomoku-Online/
├── backend/                    # 后端代码
│   ├── main.go                # 主程序入口
│   ├── go.mod                 # Go模块依赖
│   └── internal/              # 内部包
│       ├── model/             # 数据模型
│       │   └── board.go       # 棋盘和游戏状态
│       ├── service/           # 业务逻辑
│       │   ├── ai_service.go  # AI算法服务
│       │   └── game_service.go # 游戏管理服务
│       └── controller/        # 控制器
│           ├── ai_controller.go   # AI接口控制器
│           └── game_controller.go # 游戏控制器
├── frontend/                  # 前端代码
│   ├── index.html            # HTML入口
│   ├── package.json          # NPM依赖配置
│   ├── vite.config.ts        # Vite配置
│   ├── tsconfig.json         # TypeScript配置
│   └── src/                  # 源代码
│       ├── main.ts           # 应用入口
│       ├── App.vue           # 主组件
│       ├── components/       # Vue组件
│       │   ├── Board.vue     # 棋盘组件
│       │   └── ControlPanel.vue # 控制面板
│       ├── types/            # 类型定义
│       │   └── game.ts       # 游戏类型
│       ├── services/         # API服务
│       │   └── api.ts        # HTTP请求封装
│       └── utils/            # 工具函数
│           └── gameLogic.ts  # 游戏逻辑
└── .trae/                    # 项目文档
    └── documents/            # 需求和架构文档
```

## 快速开始

### 环境要求

- **Node.js**: 18.0+ (前端)
- **Golang**: 1.21+ (后端)
- **操作系统**: Windows/macOS/Linux

### 安装依赖

#### 前端依赖
```bash
cd frontend
npm install
```

#### 后端依赖
```bash
cd backend
go mod tidy
```

### 启动项目

#### 1. 启动前端 (端口: 5173)
```bash
cd frontend
npm run dev
```

前端将在 http://localhost:5173 启动

#### 2. 启动后端 (端口: 8080)
```bash
cd backend
go run main.go
```

后端API将在 http://localhost:8080 启动

### 访问游戏

打开浏览器访问: http://localhost:5173

## API接口

### AI对战接口

#### 获取AI移动
```http
POST /api/ai/move
Content-Type: application/json

{
  "board": [[0,0,0,...], [0,1,0,...], ...],
  "currentPlayer": 2
}
```

#### 获取游戏状态
```http
GET /api/ai/status
```

#### 重置游戏
```http
POST /api/ai/reset
```

### 在线匹配接口 (预留)

#### 开始匹配
```http
POST /api/match/start
Content-Type: application/json

{
  "playerId": "player123",
  "gameMode": "pvp"
}
```

#### 加入房间
```http
POST /api/match/join
Content-Type: application/json

{
  "roomId": "room123",
  "playerId": "player456"
}
```

## 游戏规则

1. **棋盘**: 15x15格子
2. **获胜条件**: 横、竖、斜任意方向连成5子
3. **回合制**: 玩家先手，AI后手
4. **AI策略**: 
   - 优先级1: 如果能赢，立即获胜
   - 优先级2: 阻止对手连成四子
   - 优先级3: 创造自己的连三威胁
   - 优先级4: 阻止对手连三
   - 优先级5: 在重要位置落子
   - 优先级6: 随机选择有效位置

## AI算法详解

### 启发式评估
- **连五**: 10000分 (必胜)
- **活四**: 1000分 (下一步必胜)
- **冲四**: 100分 (需要防守)
- **活三**: 100分 (威胁较大)
- **眠三**: 10分 (潜在威胁)
- **活二**: 10分 (发展空间)
- **眠二**: 1分 (基础分值)

### 搜索策略
1. 优先搜索已有棋子周围的位置
2. 评估每个位置的攻防价值
3. 选择综合评分最高的位置

## 开发指南

### 添加新功能

1. **后端新接口**: 在 `controller` 包中添加新的处理函数
2. **前端新组件**: 在 `src/components` 中创建Vue组件
3. **新的游戏逻辑**: 在 `utils/gameLogic.ts` 中添加工具函数

### 代码规范

- **Go**: 遵循 Go 官方代码规范
- **TypeScript**: 使用严格模式，完整类型注解
- **Vue**: 使用 Composition API 和 `<script setup>`

### 测试

#### 前端测试
```bash
cd frontend
npm run build  # 构建测试
```

#### 后端测试
```bash
cd backend
go test ./...  # 运行测试
```

## 部署

### 前端部署
```bash
cd frontend
npm run build
# 将 dist/ 目录部署到静态文件服务器
```

### 后端部署
```bash
cd backend
go build -o gomoku-server main.go
# 运行 gomoku-server 可执行文件
```

## 扩展功能 (规划中)

- 🌐 **在线匹配**: 玩家vs玩家实时对战
- 🔌 **WebSocket**: 实时通信支持
- 👥 **房间系统**: 创建和加入游戏房间
- 📱 **移动端适配**: 响应式移动端界面
- 🎯 **AI难度**: 多级别AI难度选择
- 📈 **排行榜**: 玩家积分和排名系统

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题或建议，请通过以下方式联系：
- 项目地址: [GitHub Repository]
- 邮箱: [Your Email]

---

**享受五子棋的乐趣！** 🎮