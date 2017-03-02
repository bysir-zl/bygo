package util

import (
	"sync"
	"time"
)

// 内存缓存
// 用于少量的数据, 时间为2分钟

type Data struct {
	data interface{}
	time int64
}

var dMap = map[string]Data{}
var lock sync.Mutex

func SetCache(key string, obj interface{}, exp time.Duration) {
	lock.Lock()
	defer lock.Unlock()

	dMap[key] = Data{data: obj, time: time.Now().Unix() + int64(exp)}
}

func GetCache(key string) (has bool, obj interface{}) {
	lock.Lock()
	defer lock.Unlock()

	if d, ok := dMap[key]; ok {
		if d.time <= time.Now().Unix() {
			delete(dMap, key)
			return
		}
		return true, d.data
	}
	return
}

func DeleteCache(key string) bool {
	lock.Lock()
	defer lock.Unlock()

	delete(dMap, key)
	return true
}

func init() {
	// gc
	go func() {
		for {
			time.Sleep(2 * time.Minute)
			now := time.Now().Unix()
			lock.Lock()
			for k, v := range dMap {
				if v.time <= now {
					delete(dMap, k)
				}
			}
			lock.Unlock()
		}
	}()
}
