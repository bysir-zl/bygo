package util

import (
	"math/rand"
	"time"
)

func Rand(min int, max int) (int) {
	sed := time.Now()
	time_int := sed.UnixNano()
	r := rand.New(rand.NewSource(time_int))
	res_int := r.Intn(max-min+1) + min
	return res_int
}
