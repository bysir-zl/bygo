package util

import (
	"unsafe"
	"reflect"
	"strconv"
)

// b2s converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// s2b converts string to a byte slice without memory allocation.
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func S2B(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// Hex string, high nibble first
func HexStrHead2ByteArr(hexString string) ([]byte, error) {
	lenString := len(hexString)
	if lenString % 2 == 1 {
		hexString = hexString + "0"
	}
	length := lenString / 2
	slice := make([]byte, length)
	rs := []rune(hexString)
	for i := 0; i < length; i++ {
		s := string(rs[i * 2 : i * 2 + 2])
		value, err := strconv.ParseInt(s, 16, 10)
		if err != nil {
			return nil, err
		}
		slice[i] = byte(value & 0xFF)
	}
	return slice, nil
}

func HexStrHead2String(hexString string) string {
	bs, err := HexStrHead2ByteArr(hexString)
	if err != nil {
		return ""
	}
	return B2S(bs)
}