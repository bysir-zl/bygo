package cache

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type BRedis struct {
	*redis.Pool
	prefix string
}

func NewRedis(address string) *BRedis {
	if address == "" {
		return nil
	}
	var pool = &redis.Pool{
		MaxIdle:     80,
		MaxActive:   12000, // max number of connections
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 60 {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	cache := BRedis{
		Pool: pool,
	}
	return &cache
}

// 设置一个Key前缀
func (p *BRedis) SetPrefix(prefix string) {
	p.prefix = prefix
}

func (p *BRedis) HGETALL(tableName string) (mapper map[string]interface{}, err error) {
	tableName = p.prefix + tableName
	c := p.Get()
	defer c.Close()

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

func (p *BRedis) HMSETALL(table string, mapper map[string]interface{}, expire int) error {
	table = p.prefix + table
	params := []interface{}{}
	for key, value := range mapper {
		params = append(params, key, value)
	}

	params = append([]interface{}{table}, params...)
	c := p.Get()
	defer c.Close()

	if expire != 0 {
		params = append(params, "EX", expire)
	}
	_, err := c.Do("HMSET", params...)
	if err != nil {
		return err
	}

	return nil
}

func (p *BRedis) HMGETOne(tableName string, key string) (value string, err error) {
	tableName = p.prefix + tableName
	c := p.Get()
	defer c.Close()
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

func (p *BRedis) HMSET(tableName string, key string, value interface{}, expire int) (err error) {
	tableName = p.prefix + tableName
	c := p.Get()
	defer c.Close()
	if err != nil {
		return
	}
	if expire != 0 {
		_, err = c.Do("HMSET", tableName, key, value, "EX", expire)
	} else {
		_, err = c.Do("HMSET", tableName, key, value)
	}

	return
}

func (p *BRedis) SET(key string, value interface{}, expire int) (err error) {
	key = p.prefix + key
	c := p.Get()
	defer c.Close()
	if expire != 0 {
		_, err = c.Do("SET", key, value, "EX", expire)
	} else {
		_, err = c.Do("SET", key, value)
	}

	return
}

func (p *BRedis) GET(key string) (str string, err error) {
	key = p.prefix + key
	c := p.Get()

	defer c.Close()

	value, err := c.Do("GET", key)
	if err != nil {
		return
	}
	if value != nil {
		str = string(value.([]uint8))
	}
	return
}

func (p *BRedis) RPUSH(key string, value interface{}) (err error) {
	key = p.prefix + key
	c := p.Get()
	defer c.Close()

	_, err = c.Do("RPUSH", key, value)
	return
}

func (p *BRedis) RPOP(key string, value interface{}) (data interface{}, err error) {
	key = p.prefix + key
	c := p.Get()
	defer c.Close()
	data, err = c.Do("RPOP", key, value)
	return
}

func (p *BRedis) LRANGE(key string, start, end int) (data []interface{}, err error) {
	key = p.prefix + key
	c := p.Get()
	defer c.Close()
	reply, err := c.Do("LRANGE", key, start, end)
	if err != nil {
		return
	}
	data, _ = reply.([]interface{})

	return
}

func (p *BRedis) DEL(key string) (err error) {
	key = p.prefix + key
	c := p.Get()
	defer c.Close()

	_, err = c.Do("DEL", key)
	if err != nil {
		return
	}
	return
}

func (p *BRedis) HDEL(tableName string, keys ...string) (err error) {
	tableName = p.prefix + tableName
	c := p.Get()
	defer c.Close()

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
func (p *BRedis) Lock(key string) (err error) {
	key = p.prefix + key
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
func (p *BRedis) UnLock(key string) (err error) {
	key = p.prefix + key
	err = p.DEL(key)
	return
}

// 监听事件
/*
case redis.Message:
            fmt.Printf("Message: %s %s\n", n.Channel, n.Data)
        case redis.PMessage:
            fmt.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
        case redis.Subscription:
*/

func (p *BRedis) Subscribe(key string) (event chan interface{}, err error) {
	c := p.Get()
	psc := redis.PubSubConn{Conn: c}
	err = psc.Subscribe(key)
	if err != nil {
		c.Close()
		return
	}
	event = make(chan interface{}, 256)
	go func() {
		defer c.Close()
		for {
			r := psc.Receive()
			switch n := r.(type) {
			case redis.Message, redis.PMessage, redis.Subscription:
				event <- n
			case error:
				event <- n
				close(event)
				return
			}
		}
	}()
	return
}
