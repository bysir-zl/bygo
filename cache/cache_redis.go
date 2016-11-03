package cache

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type cacheRedis struct {
	Pool *redis.Pool
}

func NewCacheRedis(redisHost string) (cache cacheRedis) {
	var pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	cache = cacheRedis{
		Pool: pool,
	}
	return
}

func (p cacheRedis) Get(key string) (value string) {
	c := p.Pool.Get()
	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		c.Close()
		return
	}
	c.Close()
	return
}

func (p cacheRedis) Set(key string, value interface{}) (ok bool) {
	c := p.Pool.Get()
	_, err := redis.String(c.Do("SET", key, value))
	if err != nil {
		log.Println(err)
	}

	c.Close()
	return err != nil
}
func (p cacheRedis) Forget(key ...string) (ok bool) {
	c := p.Pool.Get()
	count, err := redis.Int(c.Do("DEL", key[0]))
	if err != nil {
		log.Println(err)
	}
	c.Close()

	return count != 0
}