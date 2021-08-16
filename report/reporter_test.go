package report

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_can_retrieve_reported_message(t *testing.T) {
	text := "Dummy Message"
	result := reportAndReceive(func() {
		Post(text)
	})
	assert.Equal(t, text, result.Text)
}

func Test_one_message_and_multiple_receivers(t *testing.T) {
	const nbListeners = 2
	text := "Dummy Message"
	var result [nbListeners]Message
	var c[nbListeners]chan bool

	received := make(chan int, nbListeners)

	for i := 0; i < nbListeners; i++ {
		go func(i int) {
			c[i] = Subscribe(func(msg Message) {
				result[i] = msg
				received <- i
			})
		}(i)
	}

	// To make sure observers are ready to receive
	time.Sleep(1 * time.Millisecond)
	Post(text)

	for i := 0; i < nbListeners; i++ {
		iReceived := <-received
		Unsubscribe(c[iReceived])
		assert.Equal(t, text, result[iReceived].Text)
	}
}

func Test_multiple_messages_and_one_receiver(t *testing.T) {
	const nbMessages = 2
	received := make(chan Message)

	c := Subscribe(func(msg Message) {
		received <- msg
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	for i := 0; i < nbMessages; i++ {
		text := fmt.Sprintf("Dummy Message %v", i)
		Post(text)
		result := <-received
		assert.Equal(t, text, result.Text)
	}
	Unsubscribe(c)
}

func Test_report_simple_message(t *testing.T) {
	text := "Normal Message"
	result := reportAndReceive(func() {
		PostText(text)
	})
	assertMessageMatch(t, text, Normal, result)
}

func Test_report_info_message(t *testing.T) {
	text := "info Message"
	result := reportAndReceive(func() {
		PostInfo(text)
	})
	assertMessageMatch(t, text, Info, result)
}

func Test_report_title_message(t *testing.T) {
	text := "Title Message"
	result := reportAndReceive(func() {
		PostTitle(text)
	})
	assertMessageMatch(t, text, Title, result)
}

func Test_report_warning_message(t *testing.T) {
	text := "Warning Message"
	result := reportAndReceive(func() {
		PostWarning(text)
	})
	assertMessageMatch(t, text, Warning, result)
}
func Test_report_error_message(t *testing.T) {
	text := "Error Message"
	result := reportAndReceive(func() {
		PostError(text)
	})
	assertMessageMatch(t, text, Error, result)
}

func assertMessageMatch(t *testing.T, text string, msgType MessageType, msg Message) {
	assert.Equal(t, text, msg.Text)
	assert.Equal(t, msgType, msg.Type)
	assert.NotZero(t, msg.Timestamp)
}

func reportAndReceive(report func()) Message {
	var result Message
	received := make(chan bool)

	c := Subscribe(func(msg Message) {
		result = msg
		received <- true
	})

	// To make sure the observer is ready to receive
	time.Sleep(1 * time.Millisecond)
	report()
	<-received
	Unsubscribe(c)
	return result
}

