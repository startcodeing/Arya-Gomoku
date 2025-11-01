package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"gomoku-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=50"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Success      bool        `json:"success"`
	Message      string      `json:"message"`
	AccessToken  string      `json:"access_token,omitempty"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	User         *model.User `json:"user,omitempty"`
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// 验证输入
	if err := s.validateRegisterRequest(req); err != nil {
		return &AuthResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return &AuthResponse{
			Success: false,
			Message: "用户名已存在",
		}, nil
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return &AuthResponse{
			Success: false,
			Message: "邮箱已被注册",
		}, nil
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 设置昵称
	nickname := req.Nickname
	if nickname == "" {
		nickname = req.Username
	}

	// 创建用户
	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		Nickname:     nickname,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		IsActive:     true,
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 创建用户统计记录
	userStats := model.UserStatistics{
		UserID: user.ID,
	}
	if err := tx.Create(&userStats).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建用户统计失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %w", err)
	}

	return &AuthResponse{
		Success: true,
		Message: "注册成功",
		User:    &user,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	// 查找用户（支持用户名或邮箱登录）
	var user model.User
	query := s.db.Where("username = ? OR email = ?", req.Username, req.Username)
	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &AuthResponse{
				Success: false,
				Message: "用户名或密码错误",
			}, nil
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户是否活跃
	if !user.IsActive {
		return &AuthResponse{
			Success: false,
			Message: "账户已被禁用，请联系管理员",
		}, nil
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return &AuthResponse{
			Success: false,
			Message: "用户名或密码错误",
		}, nil
	}

	// 更新最后活跃时间
	if err := s.db.Model(&user).Update("last_active_at", time.Now()).Error; err != nil {
		// 记录错误但不影响登录
		fmt.Printf("更新用户活跃时间失败: %v\n", err)
	}

	return &AuthResponse{
		Success: true,
		Message: "登录成功",
		User:    &user,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("id = ? AND is_active = ?", userID, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user, nil
}

// UpdateUserProfile 更新用户资料
func (s *AuthService) UpdateUserProfile(ctx context.Context, userID string, updates map[string]interface{}) error {
	// 验证更新字段
	allowedFields := map[string]bool{
		"nickname":   true,
		"avatar_url": true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if len(filteredUpdates) == 0 {
		return fmt.Errorf("没有有效的更新字段")
	}

	filteredUpdates["updated_at"] = time.Now()

	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Updates(filteredUpdates).Error; err != nil {
		return fmt.Errorf("更新用户资料失败: %w", err)
	}

	return nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	// 获取用户
	var user model.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("原密码错误")
	}

	// 验证新密码
	if err := s.validatePassword(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := s.db.Model(&user).Updates(map[string]interface{}{
		"password_hash": string(hashedPassword),
		"updated_at":    time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// DeactivateUser 停用用户
func (s *AuthService) DeactivateUser(ctx context.Context, userID string) error {
	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("停用用户失败: %w", err)
	}

	return nil
}

// validateRegisterRequest 验证注册请求
func (s *AuthService) validateRegisterRequest(req *RegisterRequest) error {
	// 验证用户名
	if err := s.validateUsername(req.Username); err != nil {
		return err
	}

	// 验证邮箱
	if err := s.validateEmail(req.Email); err != nil {
		return err
	}

	// 验证密码
	if err := s.validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

// validateUsername 验证用户名
func (s *AuthService) validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("用户名长度必须在3-20个字符之间")
	}

	// 用户名只能包含字母、数字、下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return fmt.Errorf("用户名只能包含字母、数字和下划线")
	}

	return nil
}

// validateEmail 验证邮箱
func (s *AuthService) validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("邮箱格式不正确")
	}

	return nil
}

// validatePassword 验证密码
func (s *AuthService) validatePassword(password string) error {
	if len(password) < 8 || len(password) > 50 {
		return fmt.Errorf("密码长度必须在8-50个字符之间")
	}

	// 密码必须包含至少一个字母和一个数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasNumber {
		return fmt.Errorf("密码必须包含至少一个字母和一个数字")
	}

	return nil
}