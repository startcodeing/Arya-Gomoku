package service

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// CacheType 缓存类型
type CacheType string

const (
	CacheTypeMemory CacheType = "memory"
	CacheTypeRedis  CacheType = "redis" // 预留Redis扩展
)

// CacheConfig 缓存配置
type CacheConfig struct {
	Type     CacheType `json:"type"`
	Capacity int       `json:"capacity"`
	// Redis配置（预留）
	RedisAddr     string `json:"redis_addr,omitempty"`
	RedisPassword string `json:"redis_password,omitempty"`
	RedisDB       int    `json:"redis_db,omitempty"`
}

// CacheFactory 缓存工厂
type CacheFactory struct {
	config CacheConfig
}

// NewCacheFactory 创建缓存工厂
func NewCacheFactory(config CacheConfig) *CacheFactory {
	return &CacheFactory{
		config: config,
	}
}

// CreateCache 创建缓存实例
func (f *CacheFactory) CreateCache() CacheInterface {
	switch f.config.Type {
	case CacheTypeMemory:
		return NewMemoryCache(f.config.Capacity)
	case CacheTypeRedis:
		// 预留Redis实现
		return NewRedisCache(f.config.RedisAddr, f.config.RedisPassword, f.config.RedisDB)
	default:
		// 默认使用内存缓存
		return NewMemoryCache(f.config.Capacity)
	}
}

// RedisCache Redis缓存实现（预留）
type RedisCache struct {
	addr     string
	password string
	db       int
	stats    CacheStats
}

// NewRedisCache 创建Redis缓存实例（预留实现）
func NewRedisCache(addr, password string, db int) CacheInterface {
	return &RedisCache{
		addr:     addr,
		password: password,
		db:       db,
		stats:    CacheStats{},
	}
}

// Redis缓存接口实现（预留）
func (r *RedisCache) Get(key string) (interface{}, bool) {
	// TODO: 实现Redis GET操作
	return nil, false
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	// TODO: 实现Redis SET操作
	return nil
}

func (r *RedisCache) Delete(key string) error {
	// TODO: 实现Redis DELETE操作
	return nil
}

func (r *RedisCache) Clear() error {
	// TODO: 实现Redis CLEAR操作
	return nil
}

func (r *RedisCache) Stats() CacheStats {
	// TODO: 实现Redis统计信息获取
	return r.stats
}

func (r *RedisCache) Close() error {
	// TODO: 实现Redis连接关闭
	return nil
}

// 缓存键生成策略
const (
	CacheTTLShort  = 5 * time.Minute  // 游戏进行中
	CacheTTLMedium = 30 * time.Minute // 常见棋局
	CacheTTLLong   = 2 * time.Hour    // 开局定式
)

// GenerateCacheKey 生成缓存键
func GenerateCacheKey(board [][]int, model string) string {
	boardHash := sha256.Sum256([]byte(fmt.Sprintf("%v", board)))
	return fmt.Sprintf("llm_move:%s:%x", model, boardHash)
}

// GetCacheTTL 根据游戏阶段获取缓存TTL
func GetCacheTTL(moveCount int) time.Duration {
	switch {
	case moveCount <= 10:
		return CacheTTLLong // 开局定式，缓存时间长
	case moveCount <= 20:
		return CacheTTLMedium // 中局，中等缓存时间
	default:
		return CacheTTLShort // 残局，短缓存时间
	}
}