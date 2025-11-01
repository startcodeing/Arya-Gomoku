package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"gomoku-backend/internal/model"
)

// Claims JWT声明
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT认证中间件
type AuthMiddleware struct {
	secretKey string
	db        *gorm.DB
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(secretKey string, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		secretKey: secretKey,
		db:        db,
	}
}

// RequireAuth 需要认证的中间件
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := a.extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "认证失败：" + err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := a.validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "令牌无效：" + err.Error(),
			})
			c.Abort()
			return
		}

		// 验证用户是否存在且活跃
		var user model.User
		if err := a.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "用户不存在或已被禁用",
			})
			c.Abort()
			return
		}

		// 更新用户最后活跃时间
		go func() {
			a.db.Model(&user).Update("last_active_at", time.Now())
		}()

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("user", &user)

		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (a *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "权限不足：未找到用户角色",
			})
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, requiredRole := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "权限不足：需要 " + strings.Join(roles, " 或 ") + " 角色",
		})
		c.Abort()
	}
}

// RequireAdmin 需要管理员权限的中间件
func (a *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return a.RequireRole("admin")
}

// OptionalAuth 可选认证中间件（不强制要求认证）
func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := a.extractToken(c)
		if err != nil {
			// 没有token或token格式错误，继续执行但不设置用户信息
			c.Next()
			return
		}

		claims, err := a.validateToken(token)
		if err != nil {
			// token无效，继续执行但不设置用户信息
			c.Next()
			return
		}

		// 验证用户是否存在且活跃
		var user model.User
		if err := a.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
			// 用户不存在，继续执行但不设置用户信息
			c.Next()
			return
		}

		// 更新用户最后活跃时间
		go func() {
			a.db.Model(&user).Update("last_active_at", time.Now())
		}()

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)
		c.Set("user", &user)

		c.Next()
	}
}

// extractToken 从请求中提取token
func (a *AuthMiddleware) extractToken(c *gin.Context) (string, error) {
	// 从Authorization header中提取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	// 从Cookie中提取
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		return cookie, nil
	}

	// 从查询参数中提取
	if token := c.Query("token"); token != "" {
		return token, nil
	}

	return "", fmt.Errorf("未找到认证令牌")
}

// validateToken 验证token
func (a *AuthMiddleware) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(a.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的令牌")
}

// GenerateTokens 生成访问令牌和刷新令牌
func (a *AuthMiddleware) GenerateTokens(user *model.User) (string, string, error) {
	// 生成访问令牌
	accessClaims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "arya-gomoku",
			Subject:   user.ID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "arya-gomoku",
			Subject:   user.ID,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// RefreshToken 刷新令牌
func (a *AuthMiddleware) RefreshToken(refreshTokenString string) (string, string, error) {
	claims, err := a.validateToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	// 验证用户是否存在且活跃
	var user model.User
	if err := a.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error; err != nil {
		return "", "", fmt.Errorf("用户不存在或已被禁用")
	}

	// 生成新的令牌对
	return a.GenerateTokens(&user)
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) (*model.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("用户未认证")
	}

	currentUser, ok := user.(*model.User)
	if !ok {
		return nil, fmt.Errorf("用户信息格式错误")
	}

	return currentUser, nil
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", fmt.Errorf("用户未认证")
	}

	currentUserID, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("用户ID格式错误")
	}

	return currentUserID, nil
}

// IsAuthenticated 检查是否已认证
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// HasRole 检查是否有指定角色
func HasRole(c *gin.Context, role string) bool {
	userRole, exists := c.Get("user_role")
	if !exists {
		return false
	}

	return userRole.(string) == role
}