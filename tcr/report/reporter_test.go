package report

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_can_retrieve_reported_message(t *testing.T) {
	text := "Dummy Message"
	var result Message
	received := make(chan bool)

	go Subscribe(func(msg Message) {
		result = msg
		received <- true
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	Report(text)
	<-received
	assert.Equal(t, text, result.Text)
}

func Test_one_message_and_multiple_receivers(t *testing.T) {
	const nbListeners = 2
	text := "Dummy Message"
	var result [nbListeners]Message
	received := make(chan int, nbListeners)

	for i := 0; i < nbListeners; i++ {
		go func(i int) {
			go Subscribe(func(msg Message) {
				result[i] = msg
				received <- i
			})
		}(i)
	}

	// To make sure observers are ready to receive
	time.Sleep(1 * time.Millisecond)
	Report(text)

	for i := 0; i < nbListeners; i++ {
		iReceived := <-received
		assert.Equal(t, text, result[iReceived].Text)
	}
}

func Test_multiple_messages_and_one_receiver(t *testing.T) {
	const nbMessages = 2
	received := make(chan Message)

	go Subscribe(func(msg Message) {
		received <- msg
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < nbMessages; i++ {
		text := fmt.Sprintf("Dummy Message %v", i)
		Report(text)
		result := <-received
		assert.Equal(t, text, result.Text)
	}
}
