package report

import (
	"fmt"
	"github.com/imkira/go-observer"
	"time"
)

type MessageType int

const (
	Simple MessageType = iota
	Info
	Header
	Warning
	Error
)

type Message struct {
	Type      MessageType
	Text      string
	Timestamp time.Time
}

var (
	msgProperty = observer.NewProperty(Message{Type: Simple, Text: ""})
)

func Subscribe(onChange func(msg Message)) {
	stream := msgProperty.Observe()

	val := stream.Value().(Message)
	//fmt.Printf("initial value: %v\n", val)

	for {
		select {
		// wait for changes
		case <-stream.Changes():
			// advance to next value
			stream.Next()
			val = stream.Value().(Message)
			fmt.Printf("got new value: %v\n", val)
			onChange(val)
		}
	}
}

func Report(str string) {
	var message = Message{Simple, str, time.Now()}
	fmt.Println("Reporting message:", message)
	msgProperty.Update(message)
}
