package controller

import (
	"net/http"
	"time"

	"gomoku-backend/internal/middleware"
	"gomoku-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService *service.AuthService
	authMW      *middleware.AuthMiddleware
}

// NewAuthController 创建认证控制器
func NewAuthController(authService *service.AuthService, authMW *middleware.AuthMiddleware) *AuthController {
	return &AuthController{
		authService: authService,
		authMW:      authMW,
	}
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	response, err := c.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "注册失败: " + err.Error(),
		})
		return
	}

	if !response.Success {
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	response, err := c.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "登录失败: " + err.Error(),
		})
		return
	}

	if !response.Success {
		ctx.JSON(http.StatusUnauthorized, response)
		return
	}

	// 生成JWT令牌
	accessToken, refreshToken, err := c.authMW.GenerateTokens(response.User)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成令牌失败: " + err.Error(),
		})
		return
	}

	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	// 设置Cookie（可选）
	if req.RememberMe {
		ctx.SetCookie("refresh_token", refreshToken, int(30*24*time.Hour.Seconds()), "/", "", false, true)
	} else {
		ctx.SetCookie("refresh_token", refreshToken, int(24*time.Hour.Seconds()), "/", "", false, true)
	}

	ctx.JSON(http.StatusOK, response)
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// 清除Cookie
	ctx.SetCookie("refresh_token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

// RefreshToken 刷新令牌
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// 尝试从请求体获取refresh_token
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 如果请求体中没有，尝试从Cookie获取
		refreshToken, err := ctx.Cookie("refresh_token")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "缺少刷新令牌",
			})
			return
		}
		req.RefreshToken = refreshToken
	}

	accessToken, newRefreshToken, err := c.authMW.RefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "刷新令牌失败: " + err.Error(),
		})
		return
	}

	// 更新Cookie
	ctx.SetCookie("refresh_token", newRefreshToken, int(30*24*time.Hour.Seconds()), "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"success":       true,
		"message":       "令牌刷新成功",
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

// GetProfile 获取用户资料
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	user, err := c.authService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"user":    user,
	})
}

// UpdateProfile 更新用户资料
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	if len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "没有需要更新的字段",
		})
		return
	}

	if err := c.authService.UpdateUserProfile(ctx.Request.Context(), userID, updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
	})
}

// ChangePassword 修改密码
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	userID, err := middleware.GetCurrentUserID(ctx)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未认证",
		})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := c.authService.ChangePassword(ctx.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "密码修改成功",
	})
}

// GetUserInfo 获取用户信息（公开信息）
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户名不能为空",
		})
		return
	}

	user, err := c.authService.GetUserByUsername(ctx.Request.Context(), username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	// 只返回公开信息
	publicInfo := gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
		"created_at": user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"user":    publicInfo,
	})
}

// ValidateToken 验证令牌
func (c *AuthController) ValidateToken(ctx *gin.Context) {
	user, err := middleware.GetCurrentUser(ctx)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "令牌无效",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "令牌有效",
		"user":    user,
	})
}