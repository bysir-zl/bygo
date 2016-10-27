package cache

import "github.com/bysir-zl/bygo/config"

func NewCache(driverType string) (CacheInterface) {
    switch driverType {
    case "redis":
        return NewCacheRedis(config.BConfig.RedisHost)
        break;
    }
    return NewCacheRedis(config.BConfig.RedisHost)
}
