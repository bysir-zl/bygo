package cache

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type bRedis struct {
	*redis.Pool
}

func NewRedis(ip string) *bRedis {
	if ip == "" {
		return nil
	}
	var pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ip)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	cache := bRedis{
		Pool: pool,
	}
	return &cache
}

func (p *bRedis) HGETALL(tableName string) (mapper map[string]interface{}, err error) {
	c := p.Get()
	defer func() {
		c.Close()
	}()
	reply, err := c.Do("HGETALL", tableName)
	if err != nil {
		return
	}

	mapper = map[string]interface{}{}

	kvs := reply.([]interface{})

	for i := len(kvs) - 1; i >= 0; i = i - 2 {
		key := string(kvs[i-1].([]uint8))
		mapper[key] = kvs[i]
	}
	return
}

func (p *bRedis) MHSET(table string, mapper map[string]interface{}, expire int) error {
	params := []interface{}{}
	for key, value := range mapper {
		params = append(params, key, value)
	}

	params = append([]interface{}{table}, params...)
	c := p.Get()

	defer c.Close()

	_, err := c.Do("HMSET", params...)
	if err != nil {
		return err
	}
	if expire != 0 {
		c.Do("expire", table, expire)
	}

	return nil
}

func (p *bRedis) HMGETOne(tableName string, key string) (value string, err error) {
	c := p.Get()
	defer func() {
		c.Close()
	}()
	reply, err := c.Do("HMGET", tableName, key)
	if err != nil {
		return
	}

	ga := reply.([]interface{})[0]
	if ga != nil {
		value = string(ga.([]uint8))
	}
	return
}

func (p *bRedis) HMSET(tableName string, key string, value interface{}, expire int) (err error) {
	c := p.Get()
	defer func() {
		c.Close()
	}()
	_, err = c.Do("HMSET", tableName, key, value)
	if err != nil {
		return
	}
	if expire != 0 {
		c.Do("expire", key, expire)
	}

	return
}

func (p *bRedis) SET(key string, value interface{}, expire int) (err error) {
	c := p.Get()
	defer func() {
		c.Close()
	}()
	_, err = c.Do("SET", key, value)
	if err != nil {
		return
	}
	if expire != 0 {
		c.Do("expire", key, expire)
	}
	return
}

func (p *bRedis) GET(key string) (str string, err error) {
	c := p.Get()

	defer func() {
		c.Close()
	}()
	value, err := c.Do("GET", key)
	if err != nil {
		return
	}
	if value != nil {
		str = string(value.([]uint8))
	}
	return
}

func (p *bRedis) RPUSH(key string, value interface{}) (err error) {
	c := p.Get()

	defer func() {
		c.Close()
	}()

	_, err = c.Do("RPUSH", key, value)
	return
}

func (p *bRedis) DEL(key string) (err error) {
	c := p.Get()
	defer func() {
		c.Close()
	}()
	_, err = c.Do("DEL", key)
	if err != nil {
		return
	}
	return
}

func (p *bRedis) HDEL(tableName string, keys ...string) (err error) {
	c := p.Get()
	defer func() {
		c.Close()
	}()

	ps := make([]interface{}, len(keys)+1)
	ps[0] = tableName
	for i, v := range keys {
		ps[i+1] = v
	}

	_, err = c.Do("HDEL", ps...)
	if err != nil {
		return
	}
	return
}

// 同步锁
func (p *bRedis) Lock(key string) (err error) {
	startTime := time.Now()
	for {
		s, e := p.GET(key)
		if e != nil {
			err = e
			return
		}
		// 没有值则说明没锁
		if s == "" {
			// 上锁
			p.SET(key, 1, 10)
			return
		}

		// 有值就锁上
		// 如果一直有值 并且超时4s,则说明这个锁有问题,应该删除
		if time.Now().Sub(startTime) >= time.Second*4 {
			p.DEL(key)
			err = errors.New("deadlock")
			return
		}

		time.Sleep(time.Millisecond * 10)
	}
}

// 解锁
func (p *bRedis) UnLock(key string) (err error) {
	err = p.DEL(key)
	return
}
