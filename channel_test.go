package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func TestLearnChannel(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go func() {
		time.Sleep(2 * time.Second)
		channel <- "Hello from channel"
		fmt.Println("Message sent to channel")
	}()

	message := <-channel
	t.Log(message)
	time.Sleep(1 * time.Second)
}

// =============================

func RunGoFunctionChannel(channel chan string) {
	time.Sleep(2 * time.Second)
	channel <- "Hello from goroutine function"
	fmt.Println("Message sent to channel from goroutine function")
}

func TestChannelByParameter(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go RunGoFunctionChannel(channel)

	message := <-channel
	t.Log(message)
	time.Sleep(1 * time.Second)
}

func TestChannelOnlyIn(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go RunChannelOnlyIn(channel)

	message := <-channel
	t.Log(message)
	time.Sleep(1 * time.Second)
}

func RunChannelOnlyOut(ch <-chan string) {
	message := <-ch
	fmt.Println("Received message from only-out channel:", message)
}

func RunChannelOnlyIn(ch chan<- string) {
	time.Sleep(2 * time.Second)
	ch <- "Hello from only-in channel"
	fmt.Println("Message sent to only-in channel")
}

func TestChannelOnlyOut(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go RunChannelOnlyOut(channel)
	time.Sleep(1 * time.Second)
}

func TestChannelBothInOut(t *testing.T) {
	channel := make(chan string)
	defer close(channel)

	go RunChannelOnlyIn(channel)
	go RunChannelOnlyOut(channel)

	time.Sleep(3 * time.Second)
}

func TestChannelWithRange(t *testing.T) {
	channel := make(chan string)

	go func() {
		for i := 0; i < 5; i++ {
			channel <- fmt.Sprintf("Message %d", i)
		}
		close(channel)
	}()

	for msg := range channel {
		fmt.Println("Received:", msg)
	}
	time.Sleep(1 * time.Second)
}

func TestChannelWithSelect(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		channel1 <- "Message from channel 1"
	}()

	go func() {
		time.Sleep(1 * time.Second)
		channel2 <- "Message from channel 2"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-channel1:
			fmt.Println("Received:", msg1)
		case msg2 := <-channel2:
			fmt.Println("Received:", msg2)
		}
	}
}

func TestChannelWithDefaultSelect(t *testing.T) {
	channel1 := make(chan string)
	channel2 := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		channel1 <- "Message from channel 1"
	}()
	go func() {
		time.Sleep(1 * time.Second)
		channel2 <- "Message from channel 2"
	}()

	for i := 0; i < 3; i++ {
		select {
		case msg1 := <-channel1:
			fmt.Println("Received:", msg1)
		case msg2 := <-channel2:
			fmt.Println("Received:", msg2)
		default:
			fmt.Println("No messages received, doing other work...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
