package service

import (
	"time"
)

// CacheInterface 定义缓存操作接口，支持后续扩展到Redis
type CacheInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Stats() CacheStats
	Close() error
}

// CacheStats 缓存统计信息
type CacheStats struct {
	HitCount  int64   `json:"hit_count"`
	MissCount int64   `json:"miss_count"`
	ItemCount int     `json:"item_count"`
	HitRate   float64 `json:"hit_rate"`
}

// CacheEntry 缓存条目
type CacheEntry struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
}

// IsExpired 检查缓存条目是否过期
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}