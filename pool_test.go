package goroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() interface{} {
			// Nilai default ketika pool kosong, else nill
			return "New"
		},
	}
	pool.Put("Data 1")
	pool.Put("Data 2")
	pool.Put("Data 3")

	for i := 0; i < 10; i++ {
		go func() {
			data := pool.Get()
			fmt.Println("Get from pool:", data)
			time.Sleep(1 * time.Second)
			pool.Put(data)
		}()
	}
}
