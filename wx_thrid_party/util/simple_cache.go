package util

import (
    "sync"
    "time"
)

// 简单的缓存

var dataMap sync.Map

type Data struct {
    item      interface{}
    expiresAt int64
}

func SaveData(key string, value interface{}, expiresIn int64) {
    dataMap.Store(key, Data{
        item:      value,
        expiresAt: time.Now().Unix() + expiresIn,
    })
}

func GetData(key string) (value interface{}, ok bool) {
    d, ok := dataMap.Load(key)
    if !ok {
        return
    }
    data := d.(Data)
    if data.expiresAt+10 > time.Now().Unix() {
        return nil, false
    }

    value = data.item
    return
}
