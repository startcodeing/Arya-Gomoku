package repository

import (
	"context"
	"fmt"
	"time"

	"gomoku-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StatisticsRepository 统计仓储接口
type StatisticsRepository interface {
	// 用户统计相关
	GetUserStatistics(ctx context.Context, userID string) (*model.UserStatistics, error)
	UpdateUserStatistics(ctx context.Context, stats *model.UserStatistics) error
	CreateUserStatistics(ctx context.Context, stats *model.UserStatistics) error

	// 系统统计相关
	GetSystemStatistics(ctx context.Context) (*SystemStatistics, error)
	GetDailyStatistics(ctx context.Context, date time.Time) (*DailyStatistics, error)
	GetWeeklyStatistics(ctx context.Context, startDate time.Time) ([]*DailyStatistics, error)
	GetMonthlyStatistics(ctx context.Context, year int, month int) ([]*DailyStatistics, error)

	// 游戏统计相关
	GetGameTypeStatistics(ctx context.Context) (*GameTypeStatistics, error)
	GetDifficultyStatistics(ctx context.Context) (*DifficultyStatistics, error)
	GetWinRateStatistics(ctx context.Context) (*WinRateStatistics, error)

	// 用户行为统计
	GetActiveUserStatistics(ctx context.Context, days int) (*ActiveUserStatistics, error)
	GetUserGrowthStatistics(ctx context.Context, days int) ([]*UserGrowthData, error)
	GetUserRetentionStatistics(ctx context.Context, days int) (*UserRetentionStatistics, error)

	// 排行榜相关
	GetTopPlayersByWinRate(ctx context.Context, gameType string, limit int) ([]*PlayerRanking, error)
	GetTopPlayersByGames(ctx context.Context, gameType string, limit int) ([]*PlayerRanking, error)
	GetRecentActiveUsers(ctx context.Context, limit int) ([]*ActiveUserInfo, error)
}

// SystemStatistics 系统统计
type SystemStatistics struct {
	TotalUsers       int64     `json:"total_users"`
	ActiveUsers      int64     `json:"active_users"`
	TotalGames       int64     `json:"total_games"`
	CompletedGames   int64     `json:"completed_games"`
	OngoingGames     int64     `json:"ongoing_games"`
	TotalGameTime    int64     `json:"total_game_time"`
	AvgGameDuration  float64   `json:"avg_game_duration"`
	NewUsersToday    int64     `json:"new_users_today"`
	GamesToday       int64     `json:"games_today"`
	LastUpdated      time.Time `json:"last_updated"`
}

// DailyStatistics 每日统计
type DailyStatistics struct {
	Date            time.Time `json:"date"`
	NewUsers        int64     `json:"new_users"`
	ActiveUsers     int64     `json:"active_users"`
	TotalGames      int64     `json:"total_games"`
	AIGames         int64     `json:"ai_games"`
	LLMGames        int64     `json:"llm_games"`
	PVPGames        int64     `json:"pvp_games"`
	CompletedGames  int64     `json:"completed_games"`
	AvgGameDuration float64   `json:"avg_game_duration"`
}

// GameTypeStatistics 游戏类型统计
type GameTypeStatistics struct {
	AIGames  GameTypeData `json:"ai_games"`
	LLMGames GameTypeData `json:"llm_games"`
	PVPGames GameTypeData `json:"pvp_games"`
}

// GameTypeData 游戏类型数据
type GameTypeData struct {
	Total       int64   `json:"total"`
	Completed   int64   `json:"completed"`
	Ongoing     int64   `json:"ongoing"`
	Abandoned   int64   `json:"abandoned"`
	AvgDuration float64 `json:"avg_duration"`
	WinRate     float64 `json:"win_rate"`
}

// DifficultyStatistics 难度统计
type DifficultyStatistics struct {
	Easy   DifficultyData `json:"easy"`
	Medium DifficultyData `json:"medium"`
	Hard   DifficultyData `json:"hard"`
}

// DifficultyData 难度数据
type DifficultyData struct {
	Total       int64   `json:"total"`
	UserWins    int64   `json:"user_wins"`
	AIWins      int64   `json:"ai_wins"`
	Draws       int64   `json:"draws"`
	WinRate     float64 `json:"win_rate"`
	AvgDuration float64 `json:"avg_duration"`
}

// WinRateStatistics 胜率统计
type WinRateStatistics struct {
	Overall WinRateData `json:"overall"`
	AI      WinRateData `json:"ai"`
	LLM     WinRateData `json:"llm"`
	PVP     WinRateData `json:"pvp"`
}

// WinRateData 胜率数据
type WinRateData struct {
	TotalGames int64   `json:"total_games"`
	Wins       int64   `json:"wins"`
	Losses     int64   `json:"losses"`
	Draws      int64   `json:"draws"`
	WinRate    float64 `json:"win_rate"`
}

// ActiveUserStatistics 活跃用户统计
type ActiveUserStatistics struct {
	TotalUsers      int64 `json:"total_users"`
	DailyActive     int64 `json:"daily_active"`
	WeeklyActive    int64 `json:"weekly_active"`
	MonthlyActive   int64 `json:"monthly_active"`
	DailyActiveRate float64 `json:"daily_active_rate"`
	WeeklyActiveRate float64 `json:"weekly_active_rate"`
	MonthlyActiveRate float64 `json:"monthly_active_rate"`
}

// UserGrowthData 用户增长数据
type UserGrowthData struct {
	Date      time.Time `json:"date"`
	NewUsers  int64     `json:"new_users"`
	TotalUsers int64    `json:"total_users"`
	GrowthRate float64  `json:"growth_rate"`
}

// UserRetentionStatistics 用户留存统计
type UserRetentionStatistics struct {
	Day1Retention  float64 `json:"day1_retention"`
	Day7Retention  float64 `json:"day7_retention"`
	Day30Retention float64 `json:"day30_retention"`
	AvgSessionTime float64 `json:"avg_session_time"`
	ReturnUserRate float64 `json:"return_user_rate"`
}

// ActiveUserInfo 活跃用户信息
type ActiveUserInfo struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	LastActive   time.Time `json:"last_active"`
	TotalGames   int64     `json:"total_games"`
	RecentGames  int64     `json:"recent_games"`
	WinRate      float64   `json:"win_rate"`
}

// statisticsRepository 统计仓储实现
type statisticsRepository struct {
	db *gorm.DB
}

// NewStatisticsRepository 创建统计仓储
func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{
		db: db,
	}
}

// GetUserStatistics 获取用户统计
func (r *statisticsRepository) GetUserStatistics(ctx context.Context, userID string) (*model.UserStatistics, error) {
	var stats model.UserStatistics
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户统计不存在")
		}
		return nil, fmt.Errorf("查询用户统计失败: %w", err)
	}
	return &stats, nil
}

// UpdateUserStatistics 更新用户统计
func (r *statisticsRepository) UpdateUserStatistics(ctx context.Context, stats *model.UserStatistics) error {
	if err := r.db.WithContext(ctx).Save(stats).Error; err != nil {
		return fmt.Errorf("更新用户统计失败: %w", err)
	}
	return nil
}

// CreateUserStatistics 创建用户统计
func (r *statisticsRepository) CreateUserStatistics(ctx context.Context, stats *model.UserStatistics) error {
	if err := r.db.WithContext(ctx).Create(stats).Error; err != nil {
		return fmt.Errorf("创建用户统计失败: %w", err)
	}
	return nil
}

// GetSystemStatistics 获取系统统计
func (r *statisticsRepository) GetSystemStatistics(ctx context.Context) (*SystemStatistics, error) {
	stats := &SystemStatistics{
		LastUpdated: time.Now(),
	}

	// 统计用户数量
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("统计用户总数失败: %w", err)
	}

	// 统计活跃用户（最近7天有活动）
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_active_at > ?", sevenDaysAgo).
		Count(&stats.ActiveUsers).Error; err != nil {
		return nil, fmt.Errorf("统计活跃用户失败: %w", err)
	}

	// 统计游戏数量
	var aiGames, llmGames, pvpGames int64
	
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).Count(&aiGames).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏数量失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).Count(&llmGames).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏数量失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).Count(&pvpGames).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏数量失败: %w", err)
	}

	stats.TotalGames = aiGames + llmGames + pvpGames

	// 统计已完成游戏
	var completedAI, completedLLM, completedPVP int64
	
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&completedAI).Error; err != nil {
		return nil, fmt.Errorf("统计已完成AI游戏失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&completedLLM).Error; err != nil {
		return nil, fmt.Errorf("统计已完成LLM游戏失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&completedPVP).Error; err != nil {
		return nil, fmt.Errorf("统计已完成PVP游戏失败: %w", err)
	}

	stats.CompletedGames = completedAI + completedLLM + completedPVP

	// 统计进行中游戏
	var ongoingAI, ongoingLLM, ongoingPVP int64
	
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusInProgress).Count(&ongoingAI).Error; err != nil {
		return nil, fmt.Errorf("统计进行中AI游戏失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusInProgress).Count(&ongoingLLM).Error; err != nil {
		return nil, fmt.Errorf("统计进行中LLM游戏失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusInProgress).Count(&ongoingPVP).Error; err != nil {
		return nil, fmt.Errorf("统计进行中PVP游戏失败: %w", err)
	}

	stats.OngoingGames = ongoingAI + ongoingLLM + ongoingPVP

	// 统计总游戏时长和平均时长
	var totalDuration float64
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(duration_seconds), 0) FROM (
			SELECT duration_seconds FROM ai_games WHERE status = ? AND duration_seconds > 0
			UNION ALL
			SELECT duration_seconds FROM llm_games WHERE status = ? AND duration_seconds > 0
			UNION ALL
			SELECT duration_seconds FROM pvp_games WHERE status = ? AND duration_seconds > 0
		) AS all_durations
	`, model.GameStatusCompleted, model.GameStatusCompleted, model.GameStatusCompleted).
		Scan(&totalDuration).Error; err != nil {
		return nil, fmt.Errorf("统计总游戏时长失败: %w", err)
	}

	stats.TotalGameTime = int64(totalDuration)
	if stats.CompletedGames > 0 {
		stats.AvgGameDuration = totalDuration / float64(stats.CompletedGames)
	}

	// 统计今日新用户
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ?", today).
		Count(&stats.NewUsersToday).Error; err != nil {
		return nil, fmt.Errorf("统计今日新用户失败: %w", err)
	}

	// 统计今日游戏数量
	var todayAI, todayLLM, todayPVP int64
	
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("created_at >= ?", today).Count(&todayAI).Error; err != nil {
		return nil, fmt.Errorf("统计今日AI游戏失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("created_at >= ?", today).Count(&todayLLM).Error; err != nil {
		return nil, fmt.Errorf("统计今日LLM游戏失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("created_at >= ?", today).Count(&todayPVP).Error; err != nil {
		return nil, fmt.Errorf("统计今日PVP游戏失败: %w", err)
	}

	stats.GamesToday = todayAI + todayLLM + todayPVP

	return stats, nil
}

// GetDailyStatistics 获取每日统计
func (r *statisticsRepository) GetDailyStatistics(ctx context.Context, date time.Time) (*DailyStatistics, error) {
	stats := &DailyStatistics{
		Date: date.Truncate(24 * time.Hour),
	}

	startOfDay := stats.Date
	endOfDay := startOfDay.Add(24 * time.Hour)

	// 统计新用户
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&stats.NewUsers).Error; err != nil {
		return nil, fmt.Errorf("统计每日新用户失败: %w", err)
	}

	// 统计活跃用户
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_active_at >= ? AND last_active_at < ?", startOfDay, endOfDay).
		Count(&stats.ActiveUsers).Error; err != nil {
		return nil, fmt.Errorf("统计每日活跃用户失败: %w", err)
	}

	// 统计各类游戏
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&stats.AIGames).Error; err != nil {
		return nil, fmt.Errorf("统计每日AI游戏失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&stats.LLMGames).Error; err != nil {
		return nil, fmt.Errorf("统计每日LLM游戏失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&stats.PVPGames).Error; err != nil {
		return nil, fmt.Errorf("统计每日PVP游戏失败: %w", err)
	}

	stats.TotalGames = stats.AIGames + stats.LLMGames + stats.PVPGames

	// 统计已完成游戏和平均时长
	var avgDuration float64
	if err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COUNT(*) as completed_games,
			COALESCE(AVG(duration), 0) as avg_duration
		FROM (
			SELECT duration FROM ai_games 
			WHERE created_at >= ? AND created_at < ? AND status = ?
			UNION ALL
			SELECT duration FROM llm_games 
			WHERE created_at >= ? AND created_at < ? AND status = ?
			UNION ALL
			SELECT duration FROM pvp_games 
			WHERE created_at >= ? AND created_at < ? AND status = ?
		) AS daily_games
	`, startOfDay, endOfDay, model.GameStatusCompleted,
		startOfDay, endOfDay, model.GameStatusCompleted,
		startOfDay, endOfDay, model.GameStatusCompleted).
		Scan(&struct {
			CompletedGames int64   `gorm:"column:completed_games"`
			AvgDuration    float64 `gorm:"column:avg_duration"`
		}{CompletedGames: stats.CompletedGames, AvgDuration: avgDuration}).Error; err != nil {
		return nil, fmt.Errorf("统计每日完成游戏失败: %w", err)
	}

	stats.AvgGameDuration = avgDuration

	return stats, nil
}

// GetWeeklyStatistics 获取周统计
func (r *statisticsRepository) GetWeeklyStatistics(ctx context.Context, startDate time.Time) ([]*DailyStatistics, error) {
	var weeklyStats []*DailyStatistics

	for i := 0; i < 7; i++ {
		date := startDate.AddDate(0, 0, i)
		dailyStats, err := r.GetDailyStatistics(ctx, date)
		if err != nil {
			return nil, fmt.Errorf("获取第%d天统计失败: %w", i+1, err)
		}
		weeklyStats = append(weeklyStats, dailyStats)
	}

	return weeklyStats, nil
}

// GetMonthlyStatistics 获取月统计
func (r *statisticsRepository) GetMonthlyStatistics(ctx context.Context, year int, month int) ([]*DailyStatistics, error) {
	var monthlyStats []*DailyStatistics

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	for date := startDate; date.Before(endDate); date = date.AddDate(0, 0, 1) {
		dailyStats, err := r.GetDailyStatistics(ctx, date)
		if err != nil {
			return nil, fmt.Errorf("获取%s统计失败: %w", date.Format("2006-01-02"), err)
		}
		monthlyStats = append(monthlyStats, dailyStats)
	}

	return monthlyStats, nil
}

// GetGameTypeStatistics 获取游戏类型统计
func (r *statisticsRepository) GetGameTypeStatistics(ctx context.Context) (*GameTypeStatistics, error) {
	stats := &GameTypeStatistics{}

	// 统计AI游戏
	aiStats := &stats.AIGames
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).Count(&aiStats.Total).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏总数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&aiStats.Completed).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏完成数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusInProgress).Count(&aiStats.Ongoing).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏进行中数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusAbandoned).Count(&aiStats.Abandoned).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏放弃数失败: %w", err)
	}

	// 统计AI游戏平均时长和胜率
	var aiAvgDuration float64
	var aiWinCount int64
	
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusCompleted).
		Select("AVG(duration)").Scan(&aiAvgDuration).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏平均时长失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "user").
		Count(&aiWinCount).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏胜利数失败: %w", err)
	}

	aiStats.AvgDuration = aiAvgDuration
	if aiStats.Completed > 0 {
		aiStats.WinRate = float64(aiWinCount) / float64(aiStats.Completed) * 100
	}

	// 统计LLM游戏（类似逻辑）
	llmStats := &stats.LLMGames
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).Count(&llmStats.Total).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏总数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&llmStats.Completed).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏完成数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusInProgress).Count(&llmStats.Ongoing).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏进行中数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusAbandoned).Count(&llmStats.Abandoned).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏放弃数失败: %w", err)
	}

	var llmAvgDuration float64
	var llmWinCount int64
	
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusCompleted).
		Select("AVG(duration)").Scan(&llmAvgDuration).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏平均时长失败: %w", err)
	}
	
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "user").
		Count(&llmWinCount).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏胜利数失败: %w", err)
	}

	llmStats.AvgDuration = llmAvgDuration
	if llmStats.Completed > 0 {
		llmStats.WinRate = float64(llmWinCount) / float64(llmStats.Completed) * 100
	}

	// 统计PVP游戏（类似逻辑）
	pvpStats := &stats.PVPGames
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).Count(&pvpStats.Total).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏总数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&pvpStats.Completed).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏完成数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusInProgress).Count(&pvpStats.Ongoing).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏进行中数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusAbandoned).Count(&pvpStats.Abandoned).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏放弃数失败: %w", err)
	}

	var pvpAvgDuration float64
	
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusCompleted).
		Select("AVG(duration)").Scan(&pvpAvgDuration).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏平均时长失败: %w", err)
	}

	pvpStats.AvgDuration = pvpAvgDuration
	// PVP游戏胜率计算较复杂，这里简化处理
	pvpStats.WinRate = 50.0 // 理论上PVP游戏整体胜率应该是50%

	return stats, nil
}

// GetDifficultyStatistics 获取难度统计
func (r *statisticsRepository) GetDifficultyStatistics(ctx context.Context) (*DifficultyStatistics, error) {
	stats := &DifficultyStatistics{}

	difficulties := []struct {
		difficulty model.Difficulty
		data       *DifficultyData
	}{
		{model.DifficultyEasy, &stats.Easy},
		{model.DifficultyMedium, &stats.Medium},
		{model.DifficultyHard, &stats.Hard},
	}

	for _, d := range difficulties {
		// 统计AI游戏
		if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
			Where("difficulty = ? AND status = ?", d.difficulty, model.GameStatusCompleted).
			Count(&d.data.Total).Error; err != nil {
			return nil, fmt.Errorf("统计%s难度游戏总数失败: %w", d.difficulty, err)
		}

		if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
			Where("difficulty = ? AND status = ? AND winner = ?", 
				d.difficulty, model.GameStatusCompleted, "user").
			Count(&d.data.UserWins).Error; err != nil {
			return nil, fmt.Errorf("统计%s难度用户胜利数失败: %w", d.difficulty, err)
		}

		if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
			Where("difficulty = ? AND status = ? AND winner = ?", 
				d.difficulty, model.GameStatusCompleted, "ai").
			Count(&d.data.AIWins).Error; err != nil {
			return nil, fmt.Errorf("统计%s难度AI胜利数失败: %w", d.difficulty, err)
		}

		if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
			Where("difficulty = ? AND status = ? AND winner = ?", 
				d.difficulty, model.GameStatusCompleted, "draw").
			Count(&d.data.Draws).Error; err != nil {
			return nil, fmt.Errorf("统计%s难度平局数失败: %w", d.difficulty, err)
		}

		// 计算胜率
		if d.data.Total > 0 {
			d.data.WinRate = float64(d.data.UserWins) / float64(d.data.Total) * 100
		}

		// 统计平均时长
		if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
			Where("difficulty = ? AND status = ?", d.difficulty, model.GameStatusCompleted).
			Select("AVG(duration)").Scan(&d.data.AvgDuration).Error; err != nil {
			return nil, fmt.Errorf("统计%s难度平均时长失败: %w", d.difficulty, err)
		}
	}

	return stats, nil
}

// GetWinRateStatistics 获取胜率统计
func (r *statisticsRepository) GetWinRateStatistics(ctx context.Context) (*WinRateStatistics, error) {
	stats := &WinRateStatistics{}

	// 统计AI游戏胜率
	aiData := &stats.AI
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&aiData.TotalGames).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏总数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "user").
		Count(&aiData.Wins).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏胜利数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "ai").
		Count(&aiData.Losses).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏失败数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "draw").
		Count(&aiData.Draws).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏平局数失败: %w", err)
	}

	if aiData.TotalGames > 0 {
		aiData.WinRate = float64(aiData.Wins) / float64(aiData.TotalGames) * 100
	}

	// 统计LLM游戏胜率（类似逻辑）
	llmData := &stats.LLM
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&llmData.TotalGames).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏总数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "user").
		Count(&llmData.Wins).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏胜利数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "llm").
		Count(&llmData.Losses).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏失败数失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("status = ? AND winner = ?", model.GameStatusCompleted, "draw").
		Count(&llmData.Draws).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏平局数失败: %w", err)
	}

	if llmData.TotalGames > 0 {
		llmData.WinRate = float64(llmData.Wins) / float64(llmData.TotalGames) * 100
	}

	// 统计PVP游戏胜率（简化处理）
	pvpData := &stats.PVP
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("status = ?", model.GameStatusCompleted).Count(&pvpData.TotalGames).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏总数失败: %w", err)
	}

	// PVP游戏的胜负统计较复杂，这里简化处理
	pvpData.Wins = pvpData.TotalGames / 2
	pvpData.Losses = pvpData.TotalGames / 2
	pvpData.WinRate = 50.0

	// 统计总体胜率
	overallData := &stats.Overall
	overallData.TotalGames = aiData.TotalGames + llmData.TotalGames + pvpData.TotalGames
	overallData.Wins = aiData.Wins + llmData.Wins + pvpData.Wins
	overallData.Losses = aiData.Losses + llmData.Losses + pvpData.Losses
	overallData.Draws = aiData.Draws + llmData.Draws + pvpData.Draws

	if overallData.TotalGames > 0 {
		overallData.WinRate = float64(overallData.Wins) / float64(overallData.TotalGames) * 100
	}

	return stats, nil
}

// GetActiveUserStatistics 获取活跃用户统计
func (r *statisticsRepository) GetActiveUserStatistics(ctx context.Context, days int) (*ActiveUserStatistics, error) {
	stats := &ActiveUserStatistics{}

	// 获取总用户数
	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("获取总用户数失败: %w", err)
	}

	now := time.Now()

	// 日活跃用户
	dayAgo := now.AddDate(0, 0, -1)
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_active_at > ?", dayAgo).Count(&stats.DailyActive).Error; err != nil {
		return nil, fmt.Errorf("获取日活跃用户失败: %w", err)
	}

	// 周活跃用户
	weekAgo := now.AddDate(0, 0, -7)
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_active_at > ?", weekAgo).Count(&stats.WeeklyActive).Error; err != nil {
		return nil, fmt.Errorf("获取周活跃用户失败: %w", err)
	}

	// 月活跃用户
	monthAgo := now.AddDate(0, -1, 0)
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_active_at > ?", monthAgo).Count(&stats.MonthlyActive).Error; err != nil {
		return nil, fmt.Errorf("获取月活跃用户失败: %w", err)
	}

	// 计算活跃率
	if stats.TotalUsers > 0 {
		stats.DailyActiveRate = float64(stats.DailyActive) / float64(stats.TotalUsers) * 100
		stats.WeeklyActiveRate = float64(stats.WeeklyActive) / float64(stats.TotalUsers) * 100
		stats.MonthlyActiveRate = float64(stats.MonthlyActive) / float64(stats.TotalUsers) * 100
	}

	return stats, nil
}

// GetUserGrowthStatistics 获取用户增长统计
func (r *statisticsRepository) GetUserGrowthStatistics(ctx context.Context, days int) ([]*UserGrowthData, error) {
	var growthData []*UserGrowthData

	for i := days - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Truncate(24 * time.Hour)
		nextDate := date.Add(24 * time.Hour)

		data := &UserGrowthData{
			Date: date,
		}

		// 统计当日新用户
		if err := r.db.WithContext(ctx).Model(&model.User{}).
			Where("created_at >= ? AND created_at < ?", date, nextDate).
			Count(&data.NewUsers).Error; err != nil {
			return nil, fmt.Errorf("统计%s新用户失败: %w", date.Format("2006-01-02"), err)
		}

		// 统计截至当日的总用户数
		if err := r.db.WithContext(ctx).Model(&model.User{}).
			Where("created_at < ?", nextDate).
			Count(&data.TotalUsers).Error; err != nil {
			return nil, fmt.Errorf("统计%s总用户数失败: %w", date.Format("2006-01-02"), err)
		}

		// 计算增长率
		if i < days-1 && len(growthData) > 0 {
			prevTotal := growthData[len(growthData)-1].TotalUsers
			if prevTotal > 0 {
				data.GrowthRate = float64(data.TotalUsers-prevTotal) / float64(prevTotal) * 100
			}
		}

		growthData = append(growthData, data)
	}

	return growthData, nil
}

// GetUserRetentionStatistics 获取用户留存统计
func (r *statisticsRepository) GetUserRetentionStatistics(ctx context.Context, days int) (*UserRetentionStatistics, error) {
	stats := &UserRetentionStatistics{}

	// 这里简化实现，实际应该根据用户注册时间和后续活跃时间计算留存率
	// 由于数据模型限制，这里提供基础的统计逻辑

	now := time.Now()

	// 1日留存率（注册后第二天还活跃的用户比例）
	oneDayAgo := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)

	var newUsersYesterday, activeUsersToday int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", twoDaysAgo, oneDayAgo).
		Count(&newUsersYesterday).Error; err != nil {
		return nil, fmt.Errorf("统计昨日新用户失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ? AND last_active_at >= ?", 
			twoDaysAgo, oneDayAgo, oneDayAgo).
		Count(&activeUsersToday).Error; err != nil {
		return nil, fmt.Errorf("统计今日活跃的昨日新用户失败: %w", err)
	}

	if newUsersYesterday > 0 {
		stats.Day1Retention = float64(activeUsersToday) / float64(newUsersYesterday) * 100
	}

	// 7日留存率
	sevenDaysAgo := now.AddDate(0, 0, -7)
	eightDaysAgo := now.AddDate(0, 0, -8)

	var newUsersWeekAgo, activeUsersThisWeek int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", eightDaysAgo, sevenDaysAgo).
		Count(&newUsersWeekAgo).Error; err != nil {
		return nil, fmt.Errorf("统计一周前新用户失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ? AND last_active_at >= ?", 
			eightDaysAgo, sevenDaysAgo, sevenDaysAgo).
		Count(&activeUsersThisWeek).Error; err != nil {
		return nil, fmt.Errorf("统计本周活跃的一周前新用户失败: %w", err)
	}

	if newUsersWeekAgo > 0 {
		stats.Day7Retention = float64(activeUsersThisWeek) / float64(newUsersWeekAgo) * 100
	}

	// 30日留存率
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	thirtyOneDaysAgo := now.AddDate(0, 0, -31)

	var newUsersMonthAgo, activeUsersThisMonth int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ?", thirtyOneDaysAgo, thirtyDaysAgo).
		Count(&newUsersMonthAgo).Error; err != nil {
		return nil, fmt.Errorf("统计一月前新用户失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at >= ? AND created_at < ? AND last_active_at >= ?", 
			thirtyOneDaysAgo, thirtyDaysAgo, thirtyDaysAgo).
		Count(&activeUsersThisMonth).Error; err != nil {
		return nil, fmt.Errorf("统计本月活跃的一月前新用户失败: %w", err)
	}

	if newUsersMonthAgo > 0 {
		stats.Day30Retention = float64(activeUsersThisMonth) / float64(newUsersMonthAgo) * 100
	}

	// 平均会话时长（简化计算）
	var avgSessionTime float64
	if err := r.db.WithContext(ctx).Raw(`
		SELECT AVG(duration) FROM (
			SELECT duration FROM ai_games WHERE status = ? AND duration > 0
			UNION ALL
			SELECT duration FROM llm_games WHERE status = ? AND duration > 0
			UNION ALL
			SELECT duration FROM pvp_games WHERE status = ? AND duration > 0
		) AS all_sessions
	`, model.GameStatusCompleted, model.GameStatusCompleted, model.GameStatusCompleted).
		Scan(&avgSessionTime).Error; err != nil {
		return nil, fmt.Errorf("统计平均会话时长失败: %w", err)
	}

	stats.AvgSessionTime = avgSessionTime

	// 回归用户率（简化计算）
	var returnUsers int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("last_active_at > ? AND created_at < ?", 
			now.AddDate(0, 0, -7), now.AddDate(0, 0, -30)).
		Count(&returnUsers).Error; err != nil {
		return nil, fmt.Errorf("统计回归用户失败: %w", err)
	}

	var totalOldUsers int64
	if err := r.db.WithContext(ctx).Model(&model.User{}).
		Where("created_at < ?", now.AddDate(0, 0, -30)).
		Count(&totalOldUsers).Error; err != nil {
		return nil, fmt.Errorf("统计老用户总数失败: %w", err)
	}

	if totalOldUsers > 0 {
		stats.ReturnUserRate = float64(returnUsers) / float64(totalOldUsers) * 100
	}

	return stats, nil
}

// GetTopPlayersByWinRate 根据胜率获取排行榜
func (r *statisticsRepository) GetTopPlayersByWinRate(ctx context.Context, gameType string, limit int) ([]*PlayerRanking, error) {
	gameRepo := NewGameRepository(r.db)
	return gameRepo.GetTopPlayers(ctx, gameType, limit)
}

// GetTopPlayersByGames 根据游戏数量获取排行榜
func (r *statisticsRepository) GetTopPlayersByGames(ctx context.Context, gameType string, limit int) ([]*PlayerRanking, error) {
	var rankings []*PlayerRanking

	var query string
	switch gameType {
	case "ai":
		query = `
			SELECT 
				u.id as user_id,
				u.username,
				u.nickname,
				COUNT(*) as total_games,
				SUM(CASE WHEN winner = 'user' THEN 1 ELSE 0 END) as win_games,
				CASE 
					WHEN COUNT(*) > 0 THEN 
						ROUND(SUM(CASE WHEN winner = 'user' THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
					ELSE 0 
				END as win_rate
			FROM users u
			JOIN ai_games ag ON u.id = ag.user_id
			WHERE ag.status = ?
			GROUP BY u.id, u.username, u.nickname
			ORDER BY total_games DESC, win_rate DESC
			LIMIT ?
		`
	case "llm":
		query = `
			SELECT 
				u.id as user_id,
				u.username,
				u.nickname,
				COUNT(*) as total_games,
				SUM(CASE WHEN winner = 'user' THEN 1 ELSE 0 END) as win_games,
				CASE 
					WHEN COUNT(*) > 0 THEN 
						ROUND(SUM(CASE WHEN winner = 'user' THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
					ELSE 0 
				END as win_rate
			FROM users u
			JOIN llm_games lg ON u.id = lg.user_id
			WHERE lg.status = ?
			GROUP BY u.id, u.username, u.nickname
			ORDER BY total_games DESC, win_rate DESC
			LIMIT ?
		`
	case "pvp":
		query = `
			SELECT 
				u.id as user_id,
				u.username,
				u.nickname,
				COUNT(*) as total_games,
				SUM(CASE WHEN pg.winner_id = u.id THEN 1 ELSE 0 END) as win_games,
				CASE 
					WHEN COUNT(*) > 0 THEN 
						ROUND(SUM(CASE WHEN pg.winner_id = u.id THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
					ELSE 0 
				END as win_rate
			FROM users u
			JOIN pvp_games pg ON (u.id = pg.player1_id OR u.id = pg.player2_id)
			WHERE pg.status = ?
			GROUP BY u.id, u.username, u.nickname
			ORDER BY total_games DESC, win_rate DESC
			LIMIT ?
		`
	default:
		return nil, fmt.Errorf("不支持的游戏类型: %s", gameType)
	}

	if err := r.db.WithContext(ctx).Raw(query, model.GameStatusCompleted, limit).Scan(&rankings).Error; err != nil {
		return nil, fmt.Errorf("获取游戏数量排行榜失败: %w", err)
	}

	// 设置排名
	for i, ranking := range rankings {
		ranking.Rank = i + 1
	}

	return rankings, nil
}

// GetRecentActiveUsers 获取最近活跃用户
func (r *statisticsRepository) GetRecentActiveUsers(ctx context.Context, limit int) ([]*ActiveUserInfo, error) {
	var activeUsers []*ActiveUserInfo

	query := `
		SELECT 
			u.id as user_id,
			u.username,
			u.nickname,
			u.last_active_at,
			COALESCE(total_stats.total_games, 0) as total_games,
			COALESCE(recent_stats.recent_games, 0) as recent_games,
			COALESCE(win_stats.win_rate, 0) as win_rate
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as total_games
			FROM (
				SELECT user_id FROM ai_games WHERE status = ?
				UNION ALL
				SELECT user_id FROM llm_games WHERE status = ?
				UNION ALL
				SELECT player1_id as user_id FROM pvp_games WHERE status = ?
				UNION ALL
				SELECT player2_id as user_id FROM pvp_games WHERE status = ?
			) AS all_games
			GROUP BY user_id
		) total_stats ON u.id = total_stats.user_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as recent_games
			FROM (
				SELECT user_id FROM ai_games WHERE status = ? AND created_at > ?
				UNION ALL
				SELECT user_id FROM llm_games WHERE status = ? AND created_at > ?
				UNION ALL
				SELECT player1_id as user_id FROM pvp_games WHERE status = ? AND created_at > ?
				UNION ALL
				SELECT player2_id as user_id FROM pvp_games WHERE status = ? AND created_at > ?
			) AS recent_games
			GROUP BY user_id
		) recent_stats ON u.id = recent_stats.user_id
		LEFT JOIN (
			SELECT user_id, 
				CASE 
					WHEN COUNT(*) > 0 THEN 
						ROUND(SUM(CASE WHEN is_win = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(*), 2)
					ELSE 0 
				END as win_rate
			FROM (
				SELECT user_id, CASE WHEN winner = 'user' THEN 1 ELSE 0 END as is_win
				FROM ai_games WHERE status = ?
				UNION ALL
				SELECT user_id, CASE WHEN winner = 'user' THEN 1 ELSE 0 END as is_win
				FROM llm_games WHERE status = ?
				UNION ALL
				SELECT player1_id as user_id, CASE WHEN winner_id = player1_id THEN 1 ELSE 0 END as is_win
				FROM pvp_games WHERE status = ?
				UNION ALL
				SELECT player2_id as user_id, CASE WHEN winner_id = player2_id THEN 1 ELSE 0 END as is_win
				FROM pvp_games WHERE status = ?
			) AS win_games
			GROUP BY user_id
		) win_stats ON u.id = win_stats.user_id
		WHERE u.last_active_at IS NOT NULL
		ORDER BY u.last_active_at DESC
		LIMIT ?
	`

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	if err := r.db.WithContext(ctx).Raw(query,
		model.GameStatusCompleted, model.GameStatusCompleted, model.GameStatusCompleted, model.GameStatusCompleted,
		model.GameStatusCompleted, sevenDaysAgo, model.GameStatusCompleted, sevenDaysAgo,
		model.GameStatusCompleted, sevenDaysAgo, model.GameStatusCompleted, sevenDaysAgo,
		model.GameStatusCompleted, model.GameStatusCompleted, model.GameStatusCompleted, model.GameStatusCompleted,
		limit).Scan(&activeUsers).Error; err != nil {
		return nil, fmt.Errorf("获取最近活跃用户失败: %w", err)
	}

	return activeUsers, nil
}