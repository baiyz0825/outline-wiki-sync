// outline-wiki-sync
//
// @(#)cache.go  星期四, 六月 27, 2024
// Copyright(c) 2024, yizhuobai@Tencent. All rights reserved.

package utils

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	cacheDb *gocache.Cache
}

var XCache *Cache

func init() {
	Log.Info("初始化缓存中....")
	XCache = &Cache{cacheDb: gocache.New(time.Hour, 2*time.Hour)}
	Log.Info("初始化缓存成功，默认2h清理全局缓存")
}

// GetDataFromCache 从缓存中获取值
func (c *Cache) GetDataFromCache(key string) interface{} {
	if data, _, b := c.cacheDb.GetWithExpiration(key); b && data != nil {
		return data
	} else {
		return nil
	}
}

// SetDataToCache 设置缓存（不存在 || 已过期设置成功）
func (c *Cache) SetDataToCache(key string, data interface{}, duration time.Duration) bool {
	err := c.cacheDb.Add(key, data, duration)
	if err != nil {
		Log.Error("SetCache failure: %v", err)
		return false
	}
	return true
}

// GenCollectionCacheKey
// @Description: 生成key
// @param keyFactor
// @return string
func (c *Cache) GenCollectionCacheKey(keyFactor string) string {
	return "Collection-" + keyFactor
}
