package test

import (
	"testing"
	"time"
	"log"
)

var xChan = make(chan int, 100)

// 对于还有数据的管道，close之后，还能读
func TestCou(t *testing.T) {
	go func() {
		for i := 0; i < 5; i++ {
			xChan <- i
		}
		//time.Sleep(3*time.Second)
		close(xChan)
		log.Print(len(xChan))
		log.Print("close")
	}()

	for {
		time.Sleep(1 * time.Second)
		x, ok := <-xChan
		if !ok {
			break
		}
		log.Print(x)
	}

	t.Log("finish")
}
