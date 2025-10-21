package goroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRaceConditionWithMutex(t *testing.T) {
	var mutext sync.Mutex
	x := 0

	for i := 0; i < 1000; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				mutext.Lock()
				x = x + 1
				mutext.Unlock()
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Final value of x with mutex:", x)
}

type BackAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BackAccount) AddBalance(amount int) {
	account.RWMutex.Lock()
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BackAccount) GetBalance() int {

	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

func TestRWMutex(t *testing.T) {
	account := BackAccount{}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				account.AddBalance(1)
				fmt.Println("Current Balance:", account.GetBalance())
			}
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Final Balance:", account.GetBalance())
}
