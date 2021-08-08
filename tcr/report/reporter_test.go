package report

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_can_retrieve_reported_message(t *testing.T) {
	message := "Dummy Message"
	var result string
	received := make(chan bool)

	go Register(func(msg string) {
		result = msg
		received <- true
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	Report(message)
	<-received
	assert.Equal(t, message, result)
}

func Test_can_have_multiple_listeners(t *testing.T) {
	const nbListeners = 2
	message := "Dummy Message"
	var result [nbListeners]string
	received := make(chan int, nbListeners)

	for i := 0; i < nbListeners; i++ {
		go func(i int) {
			go Register(func(msg string) {
				result[i] = msg
				received <- i
			})
		}(i)
	}

	// To make sure observers are ready to receive
	time.Sleep(1 * time.Millisecond)
	Report(message)

	for i := 0; i < nbListeners; i++ {
		iReceived := <-received
		assert.Equal(t, message, result[iReceived])
	}
}

func Test_can_receive_multiple_messages(t *testing.T) {
	const nbMessages = 2
	received := make(chan string)

	go Register(func(msg string) {
		received <- msg
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < nbMessages; i++ {
		message := fmt.Sprintf("Dummy Message %v", i)
		Report(message)
		result := <-received
		assert.Equal(t, message, result)
	}
}
