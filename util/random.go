package util

import (
	"math/rand"
)

func Rand(min int, max int) (int) {
	//sed := time.Now()
	//time_int := sed.UnixNano()
	//r := rand.New(rand.NewSource(time_int))
	res_int := rand.Intn(max - min+1) + min
	return res_int
}
