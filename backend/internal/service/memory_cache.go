package service

import (
	"container/list"
	"sync"
	"time"
)

// MemoryCache 本地内存缓存实现
type MemoryCache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
	stats    CacheStats
	mutex    sync.RWMutex
}

// lruItem LRU缓存项
type lruItem struct {
	key   string
	entry *CacheEntry
}

// NewMemoryCache 创建新的内存缓存实例
func NewMemoryCache(capacity int) CacheInterface {
	return &MemoryCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
		stats:    CacheStats{},
	}
}

// Get 获取缓存值
func (m *MemoryCache) Get(key string) (interface{}, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	element, exists := m.items[key]
	if !exists {
		m.stats.MissCount++
		m.updateHitRate()
		return nil, false
	}

	item := element.Value.(*lruItem)
	
	// 检查是否过期
	if item.entry.IsExpired() {
		m.removeElement(element)
		m.stats.MissCount++
		m.updateHitRate()
		return nil, false
	}

	// 移动到链表头部（最近使用）
	m.order.MoveToFront(element)
	m.stats.HitCount++
	m.updateHitRate()
	
	return item.entry.Value, true
}

// Set 设置缓存值
func (m *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	entry := &CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}

	// 如果key已存在，更新值
	if element, exists := m.items[key]; exists {
		item := element.Value.(*lruItem)
		item.entry = entry
		m.order.MoveToFront(element)
		return nil
	}

	// 检查容量限制
	if len(m.items) >= m.capacity {
		m.evictLRU()
	}

	// 添加新项
	item := &lruItem{
		key:   key,
		entry: entry,
	}
	element := m.order.PushFront(item)
	m.items[key] = element
	m.stats.ItemCount = len(m.items)

	return nil
}

// Delete 删除缓存项
func (m *MemoryCache) Delete(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if element, exists := m.items[key]; exists {
		m.removeElement(element)
	}
	return nil
}

// Clear 清空所有缓存
func (m *MemoryCache) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.items = make(map[string]*list.Element)
	m.order = list.New()
	m.stats.ItemCount = 0
	return nil
}

// Stats 获取缓存统计信息
func (m *MemoryCache) Stats() CacheStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.stats
}

// Close 关闭缓存（内存缓存无需特殊关闭操作）
func (m *MemoryCache) Close() error {
	return m.Clear()
}

// evictLRU 淘汰最近最少使用的项
func (m *MemoryCache) evictLRU() {
	if m.order.Len() == 0 {
		return
	}
	
	element := m.order.Back()
	if element != nil {
		m.removeElement(element)
	}
}

// removeElement 移除指定元素
func (m *MemoryCache) removeElement(element *list.Element) {
	item := element.Value.(*lruItem)
	delete(m.items, item.key)
	m.order.Remove(element)
	m.stats.ItemCount = len(m.items)
}

// updateHitRate 更新命中率
func (m *MemoryCache) updateHitRate() {
	total := m.stats.HitCount + m.stats.MissCount
	if total > 0 {
		m.stats.HitRate = float64(m.stats.HitCount) / float64(total)
	}
}

// CleanExpired 清理过期项（可定期调用）
func (m *MemoryCache) CleanExpired() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var toRemove []*list.Element

	// 遍历所有项，找出过期的
	for element := m.order.Back(); element != nil; element = element.Prev() {
		item := element.Value.(*lruItem)
		if item.entry.IsExpired() {
			toRemove = append(toRemove, element)
		}
	}

	// 移除过期项
	for _, element := range toRemove {
		m.removeElement(element)
	}
}