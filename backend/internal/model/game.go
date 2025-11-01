package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GameStatus 游戏状态枚举
type GameStatus string

const (
	GameStatusPlaying     GameStatus = "playing"
	GameStatusInProgress  GameStatus = "in_progress"
	GameStatusCompleted   GameStatus = "completed"
	GameStatusHumanWin    GameStatus = "human_win"
	GameStatusAIWin       GameStatus = "ai_win"
	GameStatusLLMWin      GameStatus = "llm_win"
	GameStatusPlayerWin   GameStatus = "player_win"
	GameStatusDraw        GameStatus = "draw"
	GameStatusAbandoned   GameStatus = "abandoned"
)

// GameResult 游戏结果枚举
type GameResult string

const (
	GameResultWin  GameResult = "win"
	GameResultLoss GameResult = "loss"
	GameResultDraw GameResult = "draw"
)

// Difficulty 难度枚举
type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
	DifficultyExpert Difficulty = "expert"
)

// AIType AI类型枚举
type AIType string

const (
	AITypeEnhanced AIType = "enhanced"
	AITypeClassic  AIType = "classic"
)

// GameMove 游戏步骤（数据库存储版本）
type GameMove struct {
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Player    int       `json:"player"` // 1: 人类/玩家1, 2: AI/LLM/玩家2
	Timestamp time.Time `json:"timestamp"`
	ThinkTime int       `json:"think_time_ms,omitempty"` // 思考时间（毫秒）
}

// Moves 步骤数组，实现数据库存储
type Moves []GameMove

// Value 实现 driver.Valuer 接口
func (m Moves) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 实现 sql.Scanner 接口
func (m *Moves) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, m)
}

// AIGame AI对战游戏记录
type AIGame struct {
	ID           string     `json:"id" gorm:"type:varchar(36);primary_key"`
	UserID       string     `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Difficulty   Difficulty `json:"difficulty" gorm:"not null;size:20"`
	AIType       AIType     `json:"ai_type" gorm:"not null;size:20"`
	Status       GameStatus `json:"status" gorm:"not null;size:20"`
	BoardSize    int        `json:"board_size" gorm:"default:15"`
	Moves        Moves      `json:"moves" gorm:"type:text"`
	Result       *string    `json:"result" gorm:"size:20"`        // user_win, ai_win, draw, ongoing
	Winner       *string    `json:"winner" gorm:"size:20"`        // user, ai, draw
	DurationSeconds *int    `json:"duration_seconds"`
	MoveCount    int        `json:"move_count" gorm:"default:0"`
	TotalTimeMs  int        `json:"total_time_ms" gorm:"default:0"`
	UserColor    string     `json:"user_color" gorm:"default:black;size:10"`
	AIThinkingTimeMs int    `json:"ai_thinking_time_ms" gorm:"default:0"`
	GameMetadata string     `json:"game_metadata" gorm:"type:text;default:{}"`
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	EndedAt      *time.Time `json:"ended_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// LLMGame LLM对战游戏记录
type LLMGame struct {
	ID           string     `json:"id" gorm:"type:varchar(36);primary_key"`
	UserID       string     `json:"user_id" gorm:"type:varchar(36);not null;index"`
	ModelName    string     `json:"model_name" gorm:"not null;size:100"`
	Difficulty   Difficulty `json:"difficulty" gorm:"not null;size:20"`
	Status       GameStatus `json:"status" gorm:"not null;size:20"`
	BoardSize    int        `json:"board_size" gorm:"default:15"`
	Moves        Moves      `json:"moves" gorm:"type:text"`
	Result       *string    `json:"result" gorm:"size:20"`        // user_win, llm_win, draw, ongoing
	Winner       *string    `json:"winner" gorm:"size:20"`        // user, llm, draw
	DurationSeconds *int    `json:"duration_seconds"`
	MoveCount    int        `json:"move_count" gorm:"default:0"`
	TotalTimeMs  int        `json:"total_time_ms" gorm:"default:0"`
	UserColor    string     `json:"user_color" gorm:"default:black;size:10"`
	LLMResponseTimeMs int   `json:"llm_response_time_ms" gorm:"default:0"`
	GameMetadata string     `json:"game_metadata" gorm:"type:text;default:{}"`
	ConversationHistory string `json:"conversation_history" gorm:"type:text;default:[]"`
	StartedAt    *time.Time `json:"started_at"`
	FinishedAt   *time.Time `json:"finished_at"`
	EndedAt      *time.Time `json:"ended_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// PVPGame 双人对战游戏记录
type PVPGame struct {
	ID              string     `json:"id" gorm:"type:varchar(36);primary_key"`
	RoomID          string     `json:"room_id" gorm:"not null;size:50;index"`
	Player1ID       *string    `json:"player1_id" gorm:"type:varchar(36)"`
	Player2ID       *string    `json:"player2_id" gorm:"type:varchar(36)"`
	Status          GameStatus `json:"status" gorm:"not null;size:20"`
	BoardSize       int        `json:"board_size" gorm:"default:15"`
	Moves           Moves      `json:"moves" gorm:"type:text"`
	Result          *string    `json:"result" gorm:"size:20"`        // player1_win, player2_win, draw, ongoing, abandoned
	WinnerID        *string    `json:"winner_id" gorm:"type:varchar(36)"`
	DurationSeconds *int       `json:"duration_seconds"`
	MoveCount       int        `json:"move_count" gorm:"default:0"`
	CurrentPlayer   string     `json:"current_player" gorm:"default:black;size:10"`
	GameMetadata    string     `json:"game_metadata" gorm:"type:text;default:{}"`
	CurrentPlayerID *string    `json:"current_player_id" gorm:"type:varchar(36)"`
	StartedAt       *time.Time `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
	EndedAt         *time.Time `json:"ended_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	
	// 关联
	Players []PVPRoomPlayer `json:"players" gorm:"foreignKey:GameID"`
}

// PVPRoomPlayer 双人对战房间玩家
type PVPRoomPlayer struct {
	ID       string `json:"id" gorm:"type:varchar(36);primary_key"`
	GameID   string `json:"game_id" gorm:"type:varchar(36);not null;index"`
	UserID   string `json:"user_id" gorm:"type:varchar(36);not null;index"`
	PlayerNo int       `json:"player_no" gorm:"not null"` // 1 或 2
	JoinedAt time.Time `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at"`
	
	// 关联
	Game PVPGame `json:"game" gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"`
	User User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// LLMConfig LLM配置
type LLMConfig struct {
	ID        string `json:"id" gorm:"type:varchar(36);primary_key"`
	UserID    string `json:"user_id" gorm:"type:varchar(36);not null;index"`
	Name      string    `json:"name" gorm:"not null;size:100"`
	Provider  string    `json:"provider" gorm:"not null;size:50"` // openai, anthropic, etc.
	Model     string    `json:"model" gorm:"not null;size:100"`
	APIKey    string    `json:"api_key" gorm:"not null;size:500"`
	BaseURL   *string   `json:"base_url" gorm:"size:500"`
	IsDefault bool      `json:"is_default" gorm:"default:false"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate 创建前钩子
func (g *AIGame) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

func (g *LLMGame) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

func (g *PVPGame) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

func (p *PVPRoomPlayer) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (l *LLMConfig) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

// GetDuration 获取游戏时长
func (g *AIGame) GetDuration() time.Duration {
	if g.StartedAt != nil && g.EndedAt != nil {
		return g.EndedAt.Sub(*g.StartedAt)
	}
	return 0
}

func (g *LLMGame) GetDuration() time.Duration {
	if g.StartedAt != nil && g.EndedAt != nil {
		return g.EndedAt.Sub(*g.StartedAt)
	}
	return 0
}

func (g *PVPGame) GetDuration() time.Duration {
	if g.StartedAt != nil && g.EndedAt != nil {
		return g.EndedAt.Sub(*g.StartedAt)
	}
	return 0
}

// IsFinished 判断游戏是否结束
func (g *AIGame) IsFinished() bool {
	return g.Status != GameStatusPlaying
}

func (g *LLMGame) IsFinished() bool {
	return g.Status != GameStatusPlaying
}

func (g *PVPGame) IsFinished() bool {
	return g.Status != GameStatusPlaying
}