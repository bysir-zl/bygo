package util

import (
	"strconv"
	"math"
)

// 两个字串的float值是否相同
// point 表示精确到第几个小数点后
func EqualFloatString(f1, f2 string, point int) bool {
	fl1, err := strconv.ParseFloat(f1, 64)
	if err != nil {
		return false
	}
	fl2, err := strconv.ParseFloat(f2, 64)
	if err != nil {
		return false
	}

	sc := math.Pow(10, float64(point))
	return int(fl1 * sc) == int(fl2 * sc)
}