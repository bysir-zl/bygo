package cache

func NewCache(driverType string) (CacheInterface) {
    switch driverType {
    case "redis":
        return NewcacheRedis()
        break;
    }
    return NewcacheRedis()
}