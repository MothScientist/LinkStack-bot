package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

const MaxGetCacheSize int16 = 512
const CacheClearPercent float32 = 20.0 // Percentage of cache to clear when full

// GetUserCache A cache structure that is created in a single copy for the duration of the program's operation.
type GetUserCache struct {
	Cache   map[CacheCompositeKey]*GetUserCacheData
	MaxSize int16
	mu      sync.RWMutex
}

// GetUserCacheData Data for the cache and the last time (in UNIX format) of interaction with it
type GetUserCacheData struct {
	Link
	Created int64
}

// CacheCompositeKey Key for cache map
type CacheCompositeKey struct {
	TelegramId int64
	LinkId     int32
}

func (c *GetUserCache) Add(key CacheCompositeKey, data Link) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.KeyExists(key) {
		log.Printf("[WARNING][CACHE] Attempt to add an existing key")
		return
	}

	cacheSize := int16(len(c.Cache))
	if cacheSize >= c.MaxSize {
		c.delOldestCacheData()
	}

	cacheData := GetUserCacheData{
		Link:    data,
		Created: time.Now().Unix(),
	}
	c.Cache[key] = &cacheData
	cacheSize = int16(len(c.Cache))

	log.Printf("[CACHE] Add. Size: %d, occupancy: %.1f%%;", cacheSize, (float32(cacheSize)/float32(c.MaxSize))*100)
}

func (c *GetUserCache) Get(key CacheCompositeKey) *Link {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if data, exists := c.Cache[key]; exists {
		data.Created = time.Now().Unix() // Update last time usage
		cacheSize := int16(len(c.Cache))
		log.Printf("[CACHE] Get. Size: %d, occupancy: %.1f%%;", cacheSize, (float32(cacheSize)/float32(c.MaxSize))*100)
		return &data.Link
	}
	return nil
}

func (c *GetUserCache) Del(key CacheCompositeKey) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.KeyExists(key) {
		log.Print("[WARNING][CACHE] Attempt to delete non-existent key;")
		return
	}

	delete(c.Cache, key)
	cacheSize := int16(len(c.Cache))
	log.Printf("[CACHE] Del. Size: %d, occupancy: %.1f%%;", cacheSize, (float32(cacheSize)/float32(c.MaxSize))*100)
}

func (c *GetUserCache) KeyExists(key CacheCompositeKey) bool {
	// We do not block, because it is executed within the framework of one blocking (inside Add/Del)
	_, exists := c.Cache[key]
	return exists
}

func (c *GetUserCache) delOldestCacheData() {
	sizeBefore := len(c.Cache)
	occupancy_before := (float32(sizeBefore) / float32(c.MaxSize)) * 100.0

	keysToDelete := findBottomCache(c.Cache)
	fmt.Print(keysToDelete)
	for _, key := range keysToDelete {
		delete(c.Cache, key)
	}

	sizeAfter := len(c.Cache)
	log.Printf("[CACHE] Clear oldest data. Size: [before: %d, after: %d], occupancy: [before: %.1f%%, after: %.1f%%];", sizeBefore, sizeAfter, occupancy_before, (float32(sizeAfter)/float32(c.MaxSize))*100)
}

// findBottomCache Finding a cache that has fallen down the stack
func findBottomCache(cache map[CacheCompositeKey]*GetUserCacheData) []CacheCompositeKey {
	cacheSize := len(cache)
	sortedKeys := make([]CacheCompositeKey, 0, cacheSize)

	for key := range cache {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Slice(sortedKeys, func(i, j int) bool {
		return cache[sortedKeys[i]].Created > cache[sortedKeys[j]].Created // Sort in descending order
	})

	countKeysToDelete := int(float32(cacheSize) * (CacheClearPercent / 100)) // How many keys to delete
	fmt.Print(countKeysToDelete)
	return sortedKeys[:countKeysToDelete]
}

// getCacheCompositeKeyByDbData Getting a composite key for cache map
func getCacheCompositeKeyByDbData(dbData DbData) CacheCompositeKey {
	return CacheCompositeKey{
		TelegramId: dbData.TelegramId,
		LinkId:     dbData.LinkId,
	}
}

func NewGetUserCache() *GetUserCache {
	return &GetUserCache{
		Cache:   make(map[CacheCompositeKey]*GetUserCacheData),
		MaxSize: MaxGetCacheSize,
	}
}
