package test

import (
	"strconv"
	"time"
	"testing"
	"sync/atomic"
)

// 会走到 x.Y!=x.X 分支
func TestX(t *testing.T) {
	x := struct {
		X string
		Y string
	}{}

	for i := 0; i < 300000; i++ {
		go func() {
			y := strconv.FormatInt(int64(i), 10)
			x = struct {
				X string
				Y string
			}{
				X: y,
				Y: y,
			}
			if x.Y != x.X {
				t.Log("-----", x)
			}
		}()
	}

	time.Sleep(1 * time.Second)

	t.Log(x)
}

// 这个就不会
func TestY(t *testing.T) {
	v := atomic.Value{}
	for i := 0; i < 300000; i++ {
		go func() {
			y := strconv.FormatInt(int64(i), 10)

			v.Store(struct {
				X string
				Y string
			}{
				X: y,
				Y: y,
			})

			x := v.Load().(struct {
				X string
				Y string
			})
			if x.Y != x.X {
				t.Log("-----", x)
			}
		}()
	}

	time.Sleep(1 * time.Second)

	t.Log(v.Load())
}
