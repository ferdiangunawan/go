package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {

	x := 0

	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				x = x + 1
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Final value of x:", x)
}
