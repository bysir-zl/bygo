package util

import (
	"strings"
	"strconv"
)

func VersionCompareBigger(big, small string) bool {
	maxCountPoint := 0
	bigCountPoint := strings.Count(big, ".")+1
	smallCountPoint := strings.Count(small, ".")+1
	if bigCountPoint > smallCountPoint {
		maxCountPoint = bigCountPoint
	} else {
		maxCountPoint = smallCountPoint
	}

	bigNums := strings.Split(big, ".")
	smallNums := strings.Split(small, ".")

	for i := 0; i < maxCountPoint; i++ {
		bigNum := 0
		if i < len(bigNums) {
			bigNum, _ = strconv.Atoi(bigNums[i])
		}
		smallNum := 0
		if i < len(smallNums) {
			smallNum, _ = strconv.Atoi(smallNums[i])
		}
		if smallNum > bigNum {
			return false
		}

	}
	return true
}
