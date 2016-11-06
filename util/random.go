package util

import (
	"time"
	"math/rand"
)

func Rand(min int, max int) (int) {
	sed := time.Now()
	time_int := sed.UnixNano()
	r := rand.New(rand.NewSource(time_int))
	res_int := r.Intn(max - min) + min
	return res_int
}