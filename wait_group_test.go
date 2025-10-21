package goroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func RunAsyncFunction(wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	fmt.Println("Async function completed")
	time.Sleep(1 * time.Second)
}

func TestWaitGroup(t *testing.T) {
	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go RunAsyncFunction(wg)
	}

	wg.Wait()
	fmt.Println("All async functions completed")
}

func TestOne(t *testing.T) {
	wg := sync.WaitGroup{}
	once := sync.Once{}

	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			once.Do(OnlyOnce)
			wg.Done()
		}()
	}

	wg.Wait()
}

func OnlyOnce() {
	fmt.Println("Only once executed")
}
