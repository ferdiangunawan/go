package goroutine

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled:", ctx.Err())
	case <-time.After(5 * time.Second):
		fmt.Println("Completed without cancellation")
	}
}

func TestContextWithValue(t *testing.T) {
	// Membuat context dengan value
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")
	contextD := context.WithValue(contextB, "d", "D")

	fmt.Println(contextA.Value("b")) // tidak dapat, karena contextA tidak memiliki value
	fmt.Println(contextB.Value("b")) // dapat: B
	fmt.Println(contextC.Value("b")) // tidak dapat, karena contextC adalah child dari contextA
	fmt.Println(contextD.Value("b")) // dapat: B (karena contextD adalah child dari contextB)
	fmt.Println(contextD.Value("c")) // tidak dapat, karena contextD tidak memiliki akses ke contextC
	fmt.Println(contextD.Value("d")) // dapat: D
}

func TestContextGetValue(t *testing.T) {
	contextF := context.Background()
	contextF = context.WithValue(contextF, "f", "F")
	contextF = context.WithValue(contextF, "c", "C")
	contextF = context.WithValue(contextF, "b", "B")

	contextA := context.WithValue(context.Background(), "a", "A")

	fmt.Println(contextF.Value("f")) // dapat
	fmt.Println(contextF.Value("c")) // dapat milik parent
	fmt.Println(contextF.Value("b")) // tidak dapat, beda parent
	fmt.Println(contextA.Value("b")) // tidak bisa mengambil data child
}

// CreateCounter - Simulasi goroutine leak (TIDAK BAIK)
// Goroutine ini akan terus berjalan tanpa bisa dihentikan
func CreateCounter() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()

	return destination
}

// TestGoroutineLeak - Contoh goroutine leak
// Goroutine akan terus berjalan meskipun kita sudah selesai menggunakannya
func TestGoroutineLeak(t *testing.T) {
	fmt.Println("Jumlah goroutine sebelum:", runtime.NumGoroutine())

	destination := CreateCounter()
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break // kita keluar dari loop, tapi goroutine masih jalan!
		}
	}

	fmt.Println("Jumlah goroutine setelah:", runtime.NumGoroutine())
	// Goroutine masih berjalan di background - ini adalah LEAK!
}

// CreateCounterWithContext - Solusi dengan context untuk mencegah leak
// Goroutine dapat dihentikan dengan context cancellation
func CreateCounterWithContext(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				// Context dibatalkan, hentikan goroutine
				return
			case destination <- counter:
				counter++
			}
		}
	}()

	return destination
}

// TestGoroutineWithContext - Solusi untuk mencegah leak dengan context
// Goroutine akan dihentikan dengan signal dari context
func TestGoroutineWithContext(t *testing.T) {
	fmt.Println("Jumlah goroutine sebelum:", runtime.NumGoroutine())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Pastikan cancel dipanggil di semua path untuk mencegah context leak

	destination := CreateCounterWithContext(ctx)

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break // cancel akan otomatis dipanggil saat function selesai
		}
	}

	time.Sleep(100 * time.Millisecond) // Beri waktu goroutine untuk selesai
	fmt.Println("Jumlah goroutine setelah:", runtime.NumGoroutine())
	// Goroutine sudah dihentikan dengan baik - TIDAK ADA LEAK!
}
