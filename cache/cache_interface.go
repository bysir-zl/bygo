package cache

type CacheInterface interface {
	Get(key string) (value string)
	Set(key string, value interface{}) (ok bool)
	Forget(key ...string) (ok bool)
}
