package cacheSet

import (
	"fmt"
	"time"
	// "go.lwh.com/linweihao/customerComplaints/utils/err"
	t "go.lwh.com/linweihao/customerComplaints/utils/time"
	// "go.lwh.com/linweihao/customerComplaints/utils/cache"
	"go.lwh.com/linweihao/customerComplaints/factory/factoryOfCache"
)

func SetCacheOfLogOfPutSuccess(userId int) {
	// entityCache := cache.InitCache()
	entityCache := factoryOfCache.MakeEntityOfCache()
	defer entityCache.CloseCache()
	suffix := t.GetSuffix()
	cacheKey := fmt.Sprintf(
		"%s:%s:%s",
		"cc",
		"put",
		suffix)
	cacheValue := userId
	ttl := time.Second * time.Duration(60)
	entityCache.SetToCache(
		cacheKey,
		cacheValue,
		ttl)
	t.ShowTimeAndMsg("Cache set success")
	return
}
