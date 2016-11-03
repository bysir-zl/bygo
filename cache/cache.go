package cache

import (
	"github.com/bysir-zl/bygo/config"
)

func NewCache(c config.Config) CacheInterface {
	switch c.CacheDrive {
	case "redis":
		return NewCacheRedis(c.RedisHost)
		break
	}
	return NewCacheRedis(c.RedisHost)
}
