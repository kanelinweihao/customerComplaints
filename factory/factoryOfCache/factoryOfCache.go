package factoryOfCache

import (
	"go.lwh.com/linweihao/customerComplaints/utils/cache"
)

func MakeEntityOfCache() (entityCache *cache.EntityCache) {
	entityCache = cache.InitCache()
	return entityCache
}
