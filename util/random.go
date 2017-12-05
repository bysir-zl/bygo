package util

import (
	"math/rand"
	"time"
	"sync/atomic"
)

var i int64 = 0
// 随机生成int型, 并发下++
func Rand(min int, max int) (int) {
	sed := time.Now()
	timeInt := sed.UnixNano()
	timeInt += atomic.AddInt64(&i, 1)
	r := rand.New(rand.NewSource(timeInt))
	resInt := r.Intn(max-min+1) + min
	return resInt
}
