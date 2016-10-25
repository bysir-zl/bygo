package cache

func NewCache(driverType string) (CacheInterface) {
    switch driverType {
    case "redis":
        return NewCacheRedis()
        break;
    }
    return NewCacheRedis()
}