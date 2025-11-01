package repository

import (
	"context"
	"fmt"
	"time"

	"gomoku-backend/internal/model"

	"gorm.io/gorm"
)

// GameRepository 游戏仓储接口
type GameRepository interface {
	// AI游戏相关
	CreateAIGame(ctx context.Context, game *model.AIGame) error
	GetAIGameByID(ctx context.Context, id string) (*model.AIGame, error)
	UpdateAIGame(ctx context.Context, game *model.AIGame) error
	GetUserAIGames(ctx context.Context, userID string, offset, limit int) ([]*model.AIGame, int64, error)
	GetAIGamesByStatus(ctx context.Context, status model.GameStatus, offset, limit int) ([]*model.AIGame, int64, error)

	// LLM游戏相关
	CreateLLMGame(ctx context.Context, game *model.LLMGame) error
	GetLLMGameByID(ctx context.Context, id string) (*model.LLMGame, error)
	UpdateLLMGame(ctx context.Context, game *model.LLMGame) error
	GetUserLLMGames(ctx context.Context, userID string, offset, limit int) ([]*model.LLMGame, int64, error)
	GetLLMGamesByStatus(ctx context.Context, status model.GameStatus, offset, limit int) ([]*model.LLMGame, int64, error)

	// PVP游戏相关
	CreatePVPGame(ctx context.Context, game *model.PVPGame) error
	GetPVPGameByID(ctx context.Context, id string) (*model.PVPGame, error)
	UpdatePVPGame(ctx context.Context, game *model.PVPGame) error
	GetUserPVPGames(ctx context.Context, userID string, offset, limit int) ([]*model.PVPGame, int64, error)
	GetPVPGamesByStatus(ctx context.Context, status model.GameStatus, offset, limit int) ([]*model.PVPGame, int64, error)

	// 统计相关
	GetUserGameStats(ctx context.Context, userID string) (*UserGameStats, error)
	GetGameTypeStatistics(ctx context.Context) (map[string]int64, error)
	GetDifficultyStatistics(ctx context.Context) (map[string]int64, error)
	GetTopPlayers(ctx context.Context, gameType string, limit int) ([]*PlayerRanking, error)
	GetGameStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*GameStats, error)

	// 通用查询
	GetRecentGames(ctx context.Context, userID string, limit int) ([]*GameSummary, error)
	SearchGames(ctx context.Context, keyword string, offset, limit int) ([]*GameSummary, int64, error)
}

// UserGameStats 用户游戏统计
type UserGameStats struct {
	UserID       string    `json:"user_id"`
	TotalGames   int64     `json:"total_games"`
	WinCount     int64     `json:"win_count"`
	LossCount    int64     `json:"loss_count"`
	DrawCount    int64     `json:"draw_count"`
	WinRate      float64   `json:"win_rate"`
	AvgDuration  float64   `json:"avg_duration"`
	AIGames      int64     `json:"ai_games"`
	LLMGames     int64     `json:"llm_games"`
	PVPGames     int64     `json:"pvp_games"`
	LastGameTime time.Time `json:"last_game_time"`
}

// GameStats 游戏统计
type GameStats struct {
	TotalGames    int64   `json:"total_games"`
	AIGames       int64   `json:"ai_games"`
	LLMGames      int64   `json:"llm_games"`
	PVPGames      int64   `json:"pvp_games"`
	CompletedRate float64 `json:"completed_rate"`
	AvgDuration   float64 `json:"avg_duration"`
}

// PlayerRanking 玩家排名
type PlayerRanking struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Nickname string  `json:"nickname"`
	WinGames int64   `json:"win_games"`
	WinRate  float64 `json:"win_rate"`
	Rank     int     `json:"rank"`
}

// GameSummary 游戏摘要
type GameSummary struct {
	ID        string             `json:"id"`
	Type      string             `json:"type"`
	Status    model.GameStatus   `json:"status"`
	Result    model.GameResult   `json:"result"`
	UserID    string             `json:"user_id"`
	Username  string             `json:"username"`
	StartTime time.Time          `json:"start_time"`
	EndTime   *time.Time         `json:"end_time"`
	Duration  *int64             `json:"duration"`
	Moves     string             `json:"moves"` // 存储为JSON字符串
}

// gameRepository 游戏仓储实现
type gameRepository struct {
	db *gorm.DB
}

// NewGameRepository 创建游戏仓储
func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{
		db: db,
	}
}

// CreateAIGame 创建AI游戏
func (r *gameRepository) CreateAIGame(ctx context.Context, game *model.AIGame) error {
	if err := r.db.WithContext(ctx).Create(game).Error; err != nil {
		return fmt.Errorf("创建AI游戏失败: %w", err)
	}
	return nil
}

// GetAIGameByID 根据ID获取AI游戏
func (r *gameRepository) GetAIGameByID(ctx context.Context, id string) (*model.AIGame, error) {
	var game model.AIGame
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("AI游戏不存在")
		}
		return nil, fmt.Errorf("查询AI游戏失败: %w", err)
	}
	return &game, nil
}

// UpdateAIGame 更新AI游戏
func (r *gameRepository) UpdateAIGame(ctx context.Context, game *model.AIGame) error {
	if err := r.db.WithContext(ctx).Save(game).Error; err != nil {
		return fmt.Errorf("更新AI游戏失败: %w", err)
	}
	return nil
}

// GetUserAIGames 获取用户AI游戏列表
func (r *gameRepository) GetUserAIGames(ctx context.Context, userID string, offset, limit int) ([]*model.AIGame, int64, error) {
	var games []*model.AIGame
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AIGame{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户AI游戏总数失败: %w", err)
	}

	// 获取游戏列表
	if err := query.
		Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户AI游戏列表失败: %w", err)
	}

	return games, total, nil
}

// GetAIGamesByStatus 根据状态获取AI游戏列表
func (r *gameRepository) GetAIGamesByStatus(ctx context.Context, status model.GameStatus, offset, limit int) ([]*model.AIGame, int64, error) {
	var games []*model.AIGame
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AIGame{}).Where("status = ?", status)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取AI游戏总数失败: %w", err)
	}

	// 获取游戏列表
	if err := query.
		Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		return nil, 0, fmt.Errorf("获取AI游戏列表失败: %w", err)
	}

	return games, total, nil
}

// CreateLLMGame 创建LLM游戏
func (r *gameRepository) CreateLLMGame(ctx context.Context, game *model.LLMGame) error {
	if err := r.db.WithContext(ctx).Create(game).Error; err != nil {
		return fmt.Errorf("创建LLM游戏失败: %w", err)
	}
	return nil
}

// GetLLMGameByID 根据ID获取LLM游戏
func (r *gameRepository) GetLLMGameByID(ctx context.Context, id string) (*model.LLMGame, error) {
	var game model.LLMGame
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("LLM游戏不存在")
		}
		return nil, fmt.Errorf("查询LLM游戏失败: %w", err)
	}
	return &game, nil
}

// UpdateLLMGame 更新LLM游戏
func (r *gameRepository) UpdateLLMGame(ctx context.Context, game *model.LLMGame) error {
	if err := r.db.WithContext(ctx).Save(game).Error; err != nil {
		return fmt.Errorf("更新LLM游戏失败: %w", err)
	}
	return nil
}

// GetUserLLMGames 获取用户LLM游戏列表
func (r *gameRepository) GetUserLLMGames(ctx context.Context, userID string, offset, limit int) ([]*model.LLMGame, int64, error) {
	var games []*model.LLMGame
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LLMGame{}).Where("user_id = ?", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户LLM游戏总数失败: %w", err)
	}

	// 获取游戏列表
	if err := query.
		Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户LLM游戏列表失败: %w", err)
	}

	return games, total, nil
}

// GetLLMGamesByStatus 根据状态获取LLM游戏列表
func (r *gameRepository) GetLLMGamesByStatus(ctx context.Context, status model.GameStatus, offset, limit int) ([]*model.LLMGame, int64, error) {
	var games []*model.LLMGame
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LLMGame{}).Where("status = ?", status)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取LLM游戏总数失败: %w", err)
	}

	// 获取游戏列表
	if err := query.
		Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		return nil, 0, fmt.Errorf("获取LLM游戏列表失败: %w", err)
	}

	return games, total, nil
}

// CreatePVPGame 创建PVP游戏
func (r *gameRepository) CreatePVPGame(ctx context.Context, game *model.PVPGame) error {
	if err := r.db.WithContext(ctx).Create(game).Error; err != nil {
		return fmt.Errorf("创建PVP游戏失败: %w", err)
	}
	return nil
}

// GetPVPGameByID 根据ID获取PVP游戏
func (r *gameRepository) GetPVPGameByID(ctx context.Context, id string) (*model.PVPGame, error) {
	var game model.PVPGame
	if err := r.db.WithContext(ctx).
		Preload("Player1").
		Preload("Player2").
		Preload("Winner").
		Where("id = ?", id).
		First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("PVP游戏不存在")
		}
		return nil, fmt.Errorf("查询PVP游戏失败: %w", err)
	}
	return &game, nil
}

// UpdatePVPGame 更新PVP游戏
func (r *gameRepository) UpdatePVPGame(ctx context.Context, game *model.PVPGame) error {
	if err := r.db.WithContext(ctx).Save(game).Error; err != nil {
		return fmt.Errorf("更新PVP游戏失败: %w", err)
	}
	return nil
}

// GetUserPVPGames 获取用户PVP游戏列表
func (r *gameRepository) GetUserPVPGames(ctx context.Context, userID string, offset, limit int) ([]*model.PVPGame, int64, error) {
	var games []*model.PVPGame
	var total int64

	query := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("player1_id = ? OR player2_id = ?", userID, userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户PVP游戏总数失败: %w", err)
	}

	// 获取游戏列表
	if err := query.
		Preload("Player1").
		Preload("Player2").
		Preload("Winner").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户PVP游戏列表失败: %w", err)
	}

	return games, total, nil
}

// GetPVPGamesByStatus 根据状态获取PVP游戏列表
func (r *gameRepository) GetPVPGamesByStatus(ctx context.Context, status model.GameStatus, offset, limit int) ([]*model.PVPGame, int64, error) {
	var games []*model.PVPGame
	var total int64

	query := r.db.WithContext(ctx).Model(&model.PVPGame{}).Where("status = ?", status)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取PVP游戏总数失败: %w", err)
	}

	// 获取游戏列表
	if err := query.
		Preload("Player1").
		Preload("Player2").
		Preload("Winner").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&games).Error; err != nil {
		return nil, 0, fmt.Errorf("获取PVP游戏列表失败: %w", err)
	}

	return games, total, nil
}

// GetUserGameStats 获取用户游戏统计
func (r *gameRepository) GetUserGameStats(ctx context.Context, userID string) (*UserGameStats, error) {
	stats := &UserGameStats{
		UserID: userID,
	}

	// 统计AI游戏
	var aiStats struct {
		Total    int64
		Win      int64
		Lose     int64
		Draw     int64
		Duration float64
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Select("COUNT(*) as total, "+
			"SUM(CASE WHEN result = 'user_win' THEN 1 ELSE 0 END) as win, "+
			"SUM(CASE WHEN result = 'ai_win' THEN 1 ELSE 0 END) as lose, "+
			"SUM(CASE WHEN result = 'draw' THEN 1 ELSE 0 END) as draw, "+
			"AVG(duration_seconds * 1000) as duration").
		Where("user_id = ?", userID).
		Scan(&aiStats).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏失败: %w", err)
	}

	stats.AIGames = aiStats.Total
	stats.TotalGames += aiStats.Total
	stats.WinCount += aiStats.Win
	stats.LossCount += aiStats.Lose
	stats.DrawCount += aiStats.Draw

	// 统计LLM游戏
	var llmStats struct {
		Total    int64
		Win      int64
		Lose     int64
		Draw     int64
		Duration float64
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Select("COUNT(*) as total, "+
			"SUM(CASE WHEN result = 'user_win' THEN 1 ELSE 0 END) as win, "+
			"SUM(CASE WHEN result = 'llm_win' THEN 1 ELSE 0 END) as lose, "+
			"SUM(CASE WHEN result = 'draw' THEN 1 ELSE 0 END) as draw, "+
			"AVG(duration_seconds * 1000) as duration").
		Where("user_id = ?", userID).
		Scan(&llmStats).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏失败: %w", err)
	}

	stats.LLMGames = llmStats.Total
	stats.TotalGames += llmStats.Total
	stats.WinCount += llmStats.Win
	stats.LossCount += llmStats.Lose
	stats.DrawCount += llmStats.Draw

	// 统计PVP游戏
	var pvpStats struct {
		Total    int64
		Win      int64
		Lose     int64
		Draw     int64
		Duration float64
	}

	// 查询PVP游戏统计，通过pvp_room_players表关联
	if err := r.db.WithContext(ctx).Table("pvp_games").
		Joins("JOIN pvp_room_players ON pvp_games.id = pvp_room_players.game_id").
		Select("COUNT(*) as total, "+
			"SUM(CASE WHEN pvp_games.winner_id = ? THEN 1 ELSE 0 END) as win, "+
			"SUM(CASE WHEN pvp_games.winner_id != ? AND pvp_games.winner_id IS NOT NULL THEN 1 ELSE 0 END) as lose, "+
			"SUM(CASE WHEN pvp_games.winner_id IS NULL THEN 1 ELSE 0 END) as draw, "+
			"AVG(pvp_games.duration_seconds * 1000) as duration").
		Where("pvp_room_players.user_id = ?", userID, userID, userID).
		Scan(&pvpStats).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏失败: %w", err)
	}

	stats.PVPGames = pvpStats.Total
	stats.TotalGames += pvpStats.Total
	stats.WinCount += pvpStats.Win
	stats.LossCount += pvpStats.Lose
	stats.DrawCount += pvpStats.Draw

	// 计算胜率
	if stats.TotalGames > 0 {
		stats.WinRate = float64(stats.WinCount) / float64(stats.TotalGames) * 100
	}

	// 计算平均游戏时长
	totalDuration := aiStats.Duration*float64(aiStats.Total) + 
		llmStats.Duration*float64(llmStats.Total) + 
		pvpStats.Duration*float64(pvpStats.Total)
	if stats.TotalGames > 0 {
		stats.AvgDuration = totalDuration / float64(stats.TotalGames)
	}

	// 获取最后游戏时间
	var lastGameTime time.Time
	if err := r.db.WithContext(ctx).Raw(`
		SELECT MAX(created_at) FROM (
			SELECT created_at FROM ai_games WHERE user_id = ?
			UNION ALL
			SELECT created_at FROM llm_games WHERE user_id = ?
			UNION ALL
			SELECT created_at FROM pvp_games WHERE player1_id = ? OR player2_id = ?
		) AS all_games
	`, userID, userID, userID, userID).Scan(&lastGameTime).Error; err == nil {
		stats.LastGameTime = lastGameTime
	}

	return stats, nil
}

// GetGameStatsByDateRange 获取指定日期范围内的游戏统计
func (r *gameRepository) GetGameStatsByDateRange(ctx context.Context, startDate, endDate time.Time) (*GameStats, error) {
	stats := &GameStats{}

	// 统计AI游戏
	var aiCount int64
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&aiCount).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏失败: %w", err)
	}
	stats.AIGames = aiCount

	// 统计LLM游戏
	var llmCount int64
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&llmCount).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏失败: %w", err)
	}
	stats.LLMGames = llmCount

	// 统计PVP游戏
	var pvpCount int64
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&pvpCount).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏失败: %w", err)
	}
	stats.PVPGames = pvpCount

	// 总游戏数
	stats.TotalGames = aiCount + llmCount + pvpCount

	// 统计已完成游戏数
	var completedAI, completedLLM, completedPVP int64

	r.db.WithContext(ctx).Model(&model.AIGame{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, model.GameStatusCompleted).
		Count(&completedAI)

	r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, model.GameStatusCompleted).
		Count(&completedLLM)

	r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, model.GameStatusCompleted).
		Count(&completedPVP)

	completedGames := completedAI + completedLLM + completedPVP

	// 计算完成率
	if stats.TotalGames > 0 {
		stats.CompletedRate = float64(completedGames) / float64(stats.TotalGames) * 100
	}

	// 计算平均游戏时长（只统计已完成的游戏）
	var totalDuration int64
	var durationCount int64

	// AI游戏平均时长
	var aiDurationResult struct {
		AvgDuration float64 `gorm:"column:avg_duration"`
		Count       int64   `gorm:"column:count"`
	}
	r.db.WithContext(ctx).Model(&model.AIGame{}).
		Select("AVG(duration) as avg_duration, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND status = ? AND duration IS NOT NULL", startDate, endDate, model.GameStatusCompleted).
		Scan(&aiDurationResult)

	if aiDurationResult.Count > 0 {
		totalDuration += int64(aiDurationResult.AvgDuration * float64(aiDurationResult.Count))
		durationCount += aiDurationResult.Count
	}

	// LLM游戏平均时长
	var llmDurationResult struct {
		AvgDuration float64 `gorm:"column:avg_duration"`
		Count       int64   `gorm:"column:count"`
	}
	r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Select("AVG(duration) as avg_duration, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND status = ? AND duration IS NOT NULL", startDate, endDate, model.GameStatusCompleted).
		Scan(&llmDurationResult)

	if llmDurationResult.Count > 0 {
		totalDuration += int64(llmDurationResult.AvgDuration * float64(llmDurationResult.Count))
		durationCount += llmDurationResult.Count
	}

	// PVP游戏平均时长
	var pvpDurationResult struct {
		AvgDuration float64 `gorm:"column:avg_duration"`
		Count       int64   `gorm:"column:count"`
	}
	r.db.WithContext(ctx).Model(&model.PVPGame{}).
		Select("AVG(duration) as avg_duration, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND status = ? AND duration IS NOT NULL", startDate, endDate, model.GameStatusCompleted).
		Scan(&pvpDurationResult)

	if pvpDurationResult.Count > 0 {
		totalDuration += int64(pvpDurationResult.AvgDuration * float64(pvpDurationResult.Count))
		durationCount += pvpDurationResult.Count
	}

	// 计算总平均时长
	if durationCount > 0 {
		stats.AvgDuration = float64(totalDuration) / float64(durationCount)
	}

	return stats, nil
}



// GetTopPlayers 获取排行榜
func (r *gameRepository) GetTopPlayers(ctx context.Context, gameType string, limit int) ([]*PlayerRanking, error) {
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
			HAVING COUNT(*) >= 5
			ORDER BY win_rate DESC, win_games DESC
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
			HAVING COUNT(*) >= 5
			ORDER BY win_rate DESC, win_games DESC
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
			HAVING COUNT(*) >= 5
			ORDER BY win_rate DESC, win_games DESC
			LIMIT ?
		`
	default:
		return nil, fmt.Errorf("不支持的游戏类型: %s", gameType)
	}

	if err := r.db.WithContext(ctx).Raw(query, model.GameStatusCompleted, limit).Scan(&rankings).Error; err != nil {
		return nil, fmt.Errorf("获取排行榜失败: %w", err)
	}

	// 设置排名
	for i, ranking := range rankings {
		ranking.Rank = i + 1
	}

	return rankings, nil
}

// GetRecentGames 获取用户最近游戏
func (r *gameRepository) GetRecentGames(ctx context.Context, userID string, limit int) ([]*GameSummary, error) {
	var summaries []*GameSummary

	query := `
		SELECT 
			id, 'ai' as type, status, 
			CASE 
				WHEN winner = 'user' THEN 'win'
				WHEN winner = 'ai' THEN 'loss'
				WHEN winner = 'draw' THEN 'draw'
				ELSE 'unknown'
			END as result,
			user_id,
			(SELECT username FROM users WHERE id = user_id) as username,
			started_at as start_time, 
			finished_at as end_time,
			duration_seconds * 1000 as duration,
			moves
		FROM ai_games 
		WHERE user_id = ?
		UNION ALL
		SELECT 
			id, 'llm' as type, status,
			CASE 
				WHEN winner = 'user' THEN 'win'
				WHEN winner = 'llm' THEN 'loss'
				WHEN winner = 'draw' THEN 'draw'
				ELSE 'unknown'
			END as result,
			user_id,
			(SELECT username FROM users WHERE id = user_id) as username,
			started_at as start_time, 
			finished_at as end_time,
			duration_seconds * 1000 as duration,
			moves
		FROM llm_games 
		WHERE user_id = ?
		UNION ALL
		SELECT 
			id, 'pvp' as type, status, 
			CASE 
				WHEN winner_id = ? THEN 'win'
				WHEN winner_id IS NULL THEN 'draw'
				ELSE 'loss'
			END as result,
			? as user_id,
			(SELECT username FROM users WHERE id = ?) as username,
			started_at as start_time, 
			finished_at as end_time,
			duration_seconds * 1000 as duration,
			moves
		FROM pvp_games 
		WHERE player1_id = ? OR player2_id = ?
		ORDER BY start_time DESC
		LIMIT ?
	`

	if err := r.db.WithContext(ctx).Raw(query, 
		userID, userID, userID, userID, userID, userID, userID, limit).
		Scan(&summaries).Error; err != nil {
		return nil, fmt.Errorf("获取最近游戏失败: %w", err)
	}

	return summaries, nil
}

// SearchGames 搜索游戏
func (r *gameRepository) SearchGames(ctx context.Context, keyword string, offset, limit int) ([]*GameSummary, int64, error) {
	var summaries []*GameSummary
	var total int64

	// 这里简化实现，实际可以根据需要搜索用户名、游戏ID等
	query := `
		SELECT 
			id, 'ai' as type, status,
			CASE 
				WHEN winner = 'user' THEN 'win'
				WHEN winner = 'ai' THEN 'loss'
				WHEN winner = 'draw' THEN 'draw'
				ELSE 'unknown'
			END as result,
			user_id,
			(SELECT username FROM users WHERE id = user_id) as username,
			started_at as start_time, 
			finished_at as end_time,
			duration_seconds * 1000 as duration,
			moves
		FROM ai_games ag
		WHERE EXISTS (SELECT 1 FROM users u WHERE u.id = ag.user_id AND u.username ILIKE ?)
		UNION ALL
		SELECT 
			id, 'llm' as type, status,
			CASE 
				WHEN winner = 'user' THEN 'win'
				WHEN winner = 'llm' THEN 'loss'
				WHEN winner = 'draw' THEN 'draw'
				ELSE 'unknown'
			END as result,
			user_id,
			(SELECT username FROM users WHERE id = user_id) as username,
			started_at as start_time, 
			finished_at as end_time,
			duration_seconds * 1000 as duration,
			moves
		FROM llm_games lg
		WHERE EXISTS (SELECT 1 FROM users u WHERE u.id = lg.user_id AND u.username ILIKE ?)
		UNION ALL
		SELECT 
			id, 'pvp' as type, status, 
			CASE 
				WHEN winner_id IS NULL THEN 'draw'
				ELSE 'win'
			END as result,
			player1_id as user_id,
			(SELECT username FROM users WHERE id = player1_id) as username,
			started_at as start_time, 
			finished_at as end_time,
			duration_seconds * 1000 as duration,
			moves
		FROM pvp_games pg
		WHERE EXISTS (
			SELECT 1 FROM users u 
			WHERE (u.id = pg.player1_id OR u.id = pg.player2_id) 
			AND u.username ILIKE ?
		)
		ORDER BY start_time DESC
		OFFSET ? LIMIT ?
	`

	searchPattern := "%" + keyword + "%"

	if err := r.db.WithContext(ctx).Raw(query, 
		searchPattern, searchPattern, searchPattern, offset, limit).
		Scan(&summaries).Error; err != nil {
		return nil, 0, fmt.Errorf("搜索游戏失败: %w", err)
	}

	// 获取总数（简化实现）
	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT id FROM ai_games ag
			WHERE EXISTS (SELECT 1 FROM users u WHERE u.id = ag.user_id AND u.username ILIKE ?)
			UNION ALL
			SELECT id FROM llm_games lg
			WHERE EXISTS (SELECT 1 FROM users u WHERE u.id = lg.user_id AND u.username ILIKE ?)
			UNION ALL
			SELECT id FROM pvp_games pg
			WHERE EXISTS (
				SELECT 1 FROM users u 
				WHERE (u.id = pg.player1_id OR u.id = pg.player2_id) 
				AND u.username ILIKE ?
			)
		) AS total_games
	`

	if err := r.db.WithContext(ctx).Raw(countQuery, 
		searchPattern, searchPattern, searchPattern).
		Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取搜索游戏总数失败: %w", err)
	}

	return summaries, total, nil
}

// GetGameTypeStatistics 获取游戏类型统计
func (r *gameRepository) GetGameTypeStatistics(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 统计AI游戏
	var aiCount int64
	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).Count(&aiCount).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏失败: %w", err)
	}
	stats["ai"] = aiCount

	// 统计LLM游戏
	var llmCount int64
	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).Count(&llmCount).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏失败: %w", err)
	}
	stats["llm"] = llmCount

	// 统计PVP游戏
	var pvpCount int64
	if err := r.db.WithContext(ctx).Model(&model.PVPGame{}).Count(&pvpCount).Error; err != nil {
		return nil, fmt.Errorf("统计PVP游戏失败: %w", err)
	}
	stats["pvp"] = pvpCount

	return stats, nil
}

// GetDifficultyStatistics 获取难度统计
func (r *gameRepository) GetDifficultyStatistics(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 统计AI游戏难度
	var aiDifficultyStats []struct {
		Difficulty string `gorm:"column:difficulty"`
		Count      int64  `gorm:"column:count"`
	}

	if err := r.db.WithContext(ctx).Model(&model.AIGame{}).
		Select("difficulty, COUNT(*) as count").
		Group("difficulty").
		Scan(&aiDifficultyStats).Error; err != nil {
		return nil, fmt.Errorf("统计AI游戏难度失败: %w", err)
	}

	for _, stat := range aiDifficultyStats {
		stats[stat.Difficulty] += stat.Count
	}

	// 统计LLM游戏难度
	var llmDifficultyStats []struct {
		Difficulty string `gorm:"column:difficulty"`
		Count      int64  `gorm:"column:count"`
	}

	if err := r.db.WithContext(ctx).Model(&model.LLMGame{}).
		Select("difficulty, COUNT(*) as count").
		Group("difficulty").
		Scan(&llmDifficultyStats).Error; err != nil {
		return nil, fmt.Errorf("统计LLM游戏难度失败: %w", err)
	}

	for _, stat := range llmDifficultyStats {
		stats[stat.Difficulty] += stat.Count
	}

	return stats, nil
}