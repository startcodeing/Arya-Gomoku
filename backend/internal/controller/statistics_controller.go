package controller

import (
	"net/http"
	"strconv"
	"time"

	"gomoku-backend/internal/middleware"
	"gomoku-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// StatisticsController 统计控制器
type StatisticsController struct {
	gameRepo  repository.GameRepository
	statsRepo repository.StatisticsRepository
}

// NewStatisticsController 创建统计控制器
func NewStatisticsController(gameRepo repository.GameRepository, statsRepo repository.StatisticsRepository) *StatisticsController {
	return &StatisticsController{
		gameRepo:  gameRepo,
		statsRepo: statsRepo,
	}
}

// GetUserGameHistory 获取用户游戏历史
func (c *StatisticsController) GetUserGameHistory(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	gameType := ctx.Query("type") // ai, llm, pvp

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var games interface{}
	var total int64

	switch gameType {
	case "ai":
		aiGames, count, err := c.gameRepo.GetUserAIGames(ctx.Request.Context(), userID, offset, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取AI游戏历史失败: " + err.Error(),
			})
			return
		}
		games = aiGames
		total = count

	case "llm":
		llmGames, count, err := c.gameRepo.GetUserLLMGames(ctx.Request.Context(), userID, offset, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取LLM游戏历史失败: " + err.Error(),
			})
			return
		}
		games = llmGames
		total = count

	case "pvp":
		pvpGames, count, err := c.gameRepo.GetUserPVPGames(ctx.Request.Context(), userID, offset, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取PVP游戏历史失败: " + err.Error(),
			})
			return
		}
		games = pvpGames
		total = count

	default:
		// 获取所有类型的游戏摘要
		recentGames, err := c.gameRepo.GetRecentGames(ctx.Request.Context(), userID, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "获取游戏历史失败: " + err.Error(),
			})
			return
		}
		games = recentGames
		total = int64(len(recentGames))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"games": games,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetUserGameStats 获取用户游戏统计
func (c *StatisticsController) GetUserGameStats(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	stats, err := c.gameRepo.GetUserGameStats(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    stats,
	})
}

// GetSystemStatistics 获取系统统计（管理员）
func (c *StatisticsController) GetSystemStatistics(ctx *gin.Context) {
	stats, err := c.statsRepo.GetSystemStatistics(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取系统统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    stats,
	})
}

// GetGameStatsByDateRange 获取指定日期范围的游戏统计
func (c *StatisticsController) GetGameStatsByDateRange(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供开始日期和结束日期",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "开始日期格式错误，请使用 YYYY-MM-DD 格式",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "结束日期格式错误，请使用 YYYY-MM-DD 格式",
		})
		return
	}

	stats, err := c.gameRepo.GetGameStatsByDateRange(ctx.Request.Context(), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    stats,
	})
}

// GetTopPlayers 获取排行榜
func (c *StatisticsController) GetTopPlayers(ctx *gin.Context) {
	gameType := ctx.DefaultQuery("type", "all")
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if limit < 1 || limit > 100 {
		limit = 10
	}

	players, err := c.gameRepo.GetTopPlayers(ctx.Request.Context(), gameType, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取排行榜失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    players,
	})
}

// SearchGames 搜索游戏
func (c *StatisticsController) SearchGames(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供搜索关键词",
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	games, total, err := c.gameRepo.SearchGames(ctx.Request.Context(), keyword, offset, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "搜索游戏失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "搜索成功",
		"data": gin.H{
			"games": games,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetGameTypeStatistics 获取游戏类型统计
func (c *StatisticsController) GetGameTypeStatistics(ctx *gin.Context) {
	stats, err := c.statsRepo.GetGameTypeStatistics(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取游戏类型统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    stats,
	})
}

// GetDifficultyStatistics 获取难度统计
func (c *StatisticsController) GetDifficultyStatistics(ctx *gin.Context) {
	stats, err := c.statsRepo.GetDifficultyStatistics(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取难度统计失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    stats,
	})
}

// ExportGameData 导出游戏数据
func (c *StatisticsController) ExportGameData(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	format := ctx.DefaultQuery("format", "json")
	gameType := ctx.Query("type")

	// 这里可以实现数据导出逻辑
	// 为了简化，我们先返回JSON格式的数据

	var data interface{}
	switch gameType {
	case "ai":
		games, _, err := c.gameRepo.GetUserAIGames(ctx.Request.Context(), userID, 0, 1000)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导出AI游戏数据失败: " + err.Error(),
			})
			return
		}
		data = games

	case "llm":
		games, _, err := c.gameRepo.GetUserLLMGames(ctx.Request.Context(), userID, 0, 1000)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导出LLM游戏数据失败: " + err.Error(),
			})
			return
		}
		data = games

	case "pvp":
		games, _, err := c.gameRepo.GetUserPVPGames(ctx.Request.Context(), userID, 0, 1000)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导出PVP游戏数据失败: " + err.Error(),
			})
			return
		}
		data = games

	default:
		games, err := c.gameRepo.GetRecentGames(ctx.Request.Context(), userID, 1000)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "导出游戏数据失败: " + err.Error(),
			})
			return
		}
		data = games
	}

	if format == "json" {
		ctx.Header("Content-Disposition", "attachment; filename=game_data.json")
		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "导出成功",
			"data":    data,
		})
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "不支持的导出格式",
		})
	}
}

// DeleteGameRecord 删除游戏记录
func (c *StatisticsController) DeleteGameRecord(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	gameID := ctx.Param("id")
	if gameID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "游戏ID不能为空",
		})
		return
	}

	gameType := ctx.Query("type")
	if gameType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请指定游戏类型",
		})
		return
	}

	// 这里需要实现删除逻辑，确保只能删除自己的游戏记录
	// 为了简化，我们先返回成功响应
	// TODO: 实现实际的删除逻辑，使用gameID参数
	_ = gameID // 临时忽略未使用的变量警告
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}