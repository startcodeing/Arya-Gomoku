package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             string         `json:"id" gorm:"type:varchar(36);primary_key"`
	Username       string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email          string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Nickname       string         `json:"nickname" gorm:"not null;size:100"`
	PasswordHash   string         `json:"-" gorm:"not null;size:255"`
	AvatarURL      *string        `json:"avatar_url" gorm:"size:500"`
	EmailVerified  bool           `json:"email_verified" gorm:"default:false"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	Role           string         `json:"role" gorm:"default:'user';size:20"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	LastActiveAt   time.Time      `json:"last_active_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserSession 用户会话模型
type UserSession struct {
	ID           string    `json:"id" gorm:"type:varchar(36);primary_key"`
	UserID       string    `json:"user_id" gorm:"type:varchar(36);not null;index"`
	SessionToken string    `json:"session_token" gorm:"uniqueIndex;not null;size:255"`
	RefreshToken string    `json:"refresh_token" gorm:"uniqueIndex;not null;size:255"`
	IPAddress    *string   `json:"ip_address" gorm:"size:45"`
	UserAgent    *string   `json:"user_agent" gorm:"type:text"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null;index"`
	CreatedAt    time.Time `json:"created_at"`
	LastUsedAt   time.Time `json:"last_used_at"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// UserStatistics 用户统计模型
type UserStatistics struct {
	UserID              string    `json:"user_id" gorm:"type:varchar(36);primary_key"`
	TotalAIGames        int       `json:"total_ai_games" gorm:"default:0"`
	AIGamesWon          int       `json:"ai_games_won" gorm:"default:0"`
	AIGamesLost         int       `json:"ai_games_lost" gorm:"default:0"`
	AIGamesDrawn        int       `json:"ai_games_drawn" gorm:"default:0"`
	TotalLLMGames       int       `json:"total_llm_games" gorm:"default:0"`
	LLMGamesWon         int       `json:"llm_games_won" gorm:"default:0"`
	LLMGamesLost        int       `json:"llm_games_lost" gorm:"default:0"`
	LLMGamesDrawn       int       `json:"llm_games_drawn" gorm:"default:0"`
	TotalPVPGames       int       `json:"total_pvp_games" gorm:"default:0"`
	PVPGamesWon         int       `json:"pvp_games_won" gorm:"default:0"`
	PVPGamesLost        int       `json:"pvp_games_lost" gorm:"default:0"`
	PVPGamesDrawn       int       `json:"pvp_games_drawn" gorm:"default:0"`
	TotalPlayTimeMs     int64     `json:"total_play_time_ms" gorm:"default:0"`
	HighestWinStreak    int       `json:"highest_win_streak" gorm:"default:0"`
	CurrentWinStreak    int       `json:"current_win_streak" gorm:"default:0"`
	LastGameAt          *time.Time `json:"last_game_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	
	// 关联
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	u.LastActiveAt = time.Now()
	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// GetTotalGames 获取总游戏数
func (us *UserStatistics) GetTotalGames() int {
	return us.TotalAIGames + us.TotalLLMGames + us.TotalPVPGames
}

// GetTotalWins 获取总胜利数
func (us *UserStatistics) GetTotalWins() int {
	return us.AIGamesWon + us.LLMGamesWon + us.PVPGamesWon
}

// GetOverallWinRate 获取总胜率
func (us *UserStatistics) GetOverallWinRate() float64 {
	totalGames := us.GetTotalGames()
	if totalGames == 0 {
		return 0.0
	}
	return float64(us.GetTotalWins()) / float64(totalGames)
}

// GetAIWinRate 获取AI对战胜率
func (us *UserStatistics) GetAIWinRate() float64 {
	if us.TotalAIGames == 0 {
		return 0.0
	}
	return float64(us.AIGamesWon) / float64(us.TotalAIGames)
}

// GetLLMWinRate 获取LLM对战胜率
func (us *UserStatistics) GetLLMWinRate() float64 {
	if us.TotalLLMGames == 0 {
		return 0.0
	}
	return float64(us.LLMGamesWon) / float64(us.TotalLLMGames)
}

// GetPVPWinRate 获取PVP对战胜率
func (us *UserStatistics) GetPVPWinRate() float64 {
	if us.TotalPVPGames == 0 {
		return 0.0
	}
	return float64(us.PVPGamesWon) / float64(us.TotalPVPGames)
}