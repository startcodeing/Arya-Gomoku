package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gomoku-backend/internal/config"
	"gomoku-backend/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 数据库管理器
type Database struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewDatabase 创建数据库连接
func NewDatabase(cfg *config.Config) (*Database, error) {
	// 连接数据库
	var db *gorm.DB
	var err error
	
	switch cfg.Database.Type {
	case "postgres":
		db, err = connectPostgreSQL(&cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
		}
	case "sqlite":
		db, err = connectSQLite(&cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// 连接Redis (可选)
	var redisClient *redis.Client
	if cfg.Database.Type != "sqlite" {
		redisClient, err = connectRedis(&cfg.Redis)
		if err != nil {
			log.Printf("Warning: failed to connect to Redis: %v", err)
			// Redis连接失败不影响主要功能，继续运行
		}
	}

	database := &Database{
		DB:    db,
		Redis: redisClient,
	}

	return database, nil
}

// connectPostgreSQL 连接PostgreSQL数据库
func connectPostgreSQL(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 配置GORM日志
	var logLevel logger.LogLevel
	switch cfg.SSLMode {
	case "development":
		logLevel = logger.Info
	default:
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), gormConfig)
	if err != nil {
		return nil, err
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// connectSQLite 连接SQLite数据库
func connectSQLite(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// 配置GORM日志
	var logLevel logger.LogLevel
	if cfg.SSLMode == "development" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// 连接SQLite数据库
	db, err := gorm.Open(sqlite.Open(cfg.DBName), gormConfig)
	if err != nil {
		return nil, err
	}

	// 获取底层sql.DB对象进行连接池配置
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 配置连接池（SQLite通常不需要太多连接）
	sqlDB.SetMaxOpenConns(1) // SQLite只支持单个写连接
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	log.Println("Successfully connected to SQLite database")
	return db, nil
}

// connectRedis 连接Redis
func connectRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.GetRedisAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return client, nil
}

// AutoMigrate 自动迁移数据库表
func (d *Database) AutoMigrate() error {
	log.Println("Starting database migration...")
	
	// 迁移所有模型
	err := d.DB.AutoMigrate(
		&model.User{},
		&model.UserSession{},
		&model.UserStatistics{},
		&model.AIGame{},
		&model.LLMGame{},
		&model.PVPGame{},
		&model.PVPRoomPlayer{},
		&model.LLMConfig{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// CreateIndexes 创建数据库索引
func (d *Database) CreateIndexes() error {
	log.Println("Creating database indexes...")
	
	// 用户表索引
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_last_active ON users(last_active_at DESC)").Error; err != nil {
		return err
	}
	
	// AI游戏表索引
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_ai_games_created_at ON ai_games(created_at DESC)").Error; err != nil {
		return err
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_ai_games_status ON ai_games(status)").Error; err != nil {
		return err
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_ai_games_difficulty ON ai_games(difficulty)").Error; err != nil {
		return err
	}
	
	// LLM游戏表索引
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_llm_games_created_at ON llm_games(created_at DESC)").Error; err != nil {
		return err
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_llm_games_model ON llm_games(model_name)").Error; err != nil {
		return err
	}
	
	// PVP游戏表索引
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_pvp_games_created_at ON pvp_games(created_at DESC)").Error; err != nil {
		return err
	}
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_pvp_games_room_id ON pvp_games(room_id)").Error; err != nil {
		return err
	}
	
	// 会话表索引
	if err := d.DB.Exec("CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at)").Error; err != nil {
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}

// InitializeData 初始化基础数据
func (d *Database) InitializeData() error {
	log.Println("Initializing base data...")
	
	// 创建默认管理员用户（如果不存在）
	var adminUser model.User
	result := d.DB.Where("username = ?", "admin").First(&adminUser)
	if result.Error == gorm.ErrRecordNotFound {
		// 创建管理员用户
		adminUser = model.User{
			Username:      "admin",
			Email:         "admin@gomoku.com",
			Nickname:      "管理员",
			PasswordHash:  "$2a$10$nwkc.2ldZ8yzYUxnSATQEO11YRa9rK0rgYpDZ5qSJWTNcWjHsFKtC", // password123
			EmailVerified: true,
			Role:          "admin",
		}
		
		if err := d.DB.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
		
		// 创建管理员统计记录
		adminStats := model.UserStatistics{
			UserID: adminUser.ID,
		}
		if err := d.DB.Create(&adminStats).Error; err != nil {
			return fmt.Errorf("failed to create admin statistics: %w", err)
		}
		
		log.Println("Admin user created successfully")
	}

	// 创建测试用户（如果不存在）
	var testUser model.User
	result = d.DB.Where("username = ?", "testuser").First(&testUser)
	if result.Error == gorm.ErrRecordNotFound {
		// 创建测试用户
		testUser = model.User{
			Username:      "testuser",
			Email:         "test@gomoku.com",
			Nickname:      "测试用户",
			PasswordHash:  "$2a$10$nwkc.2ldZ8yzYUxnSATQEO11YRa9rK0rgYpDZ5qSJWTNcWjHsFKtC", // password123
			EmailVerified: true,
			Role:          "user",
		}
		
		if err := d.DB.Create(&testUser).Error; err != nil {
			return fmt.Errorf("failed to create test user: %w", err)
		}
		
		// 创建测试用户统计记录
		testStats := model.UserStatistics{
			UserID: testUser.ID,
		}
		if err := d.DB.Create(&testStats).Error; err != nil {
			return fmt.Errorf("failed to create test user statistics: %w", err)
		}
		
		log.Println("Test user created successfully")
	}

	log.Println("Base data initialization completed")
	return nil
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	// 关闭PostgreSQL连接
	if sqlDB, err := d.DB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing PostgreSQL connection: %v", err)
		}
	}

	// 关闭Redis连接
	if err := d.Redis.Close(); err != nil {
		log.Printf("Error closing Redis connection: %v", err)
	}

	log.Println("Database connections closed")
	return nil
}

// HealthCheck 健康检查
func (d *Database) HealthCheck(ctx context.Context) error {
	// 检查PostgreSQL
	if sqlDB, err := d.DB.DB(); err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	} else if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("PostgreSQL health check failed: %w", err)
	}

	// 检查Redis
	if err := d.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}

	return nil
}