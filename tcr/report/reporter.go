package report

import (
	"fmt"
	"github.com/imkira/go-observer"
	"time"
)

type MessageType int

const (
	Normal MessageType = iota
	Info
	Title
	Warning
	Error
)

type Message struct {
	Type      MessageType
	Text      string
	Timestamp time.Time
}

var (
	msgProperty = observer.NewProperty(Message{Type: Normal, Text: ""})
)

func Subscribe(onReport func(msg Message)) chan bool {
	stream := msgProperty.Observe()

	msg := stream.Value().(Message)
	//fmt.Printf("initial value: %v\n", msg)

	unsubscribe := make(chan bool)
	go func(s observer.Stream) {
		for {
			select {
			// wait for changes
			case <-s.Changes():
				// advance to next value
				s.Next()
				msg = s.Value().(Message)
				//fmt.Printf("got new value: %v\n", msg)
				onReport(msg)
			case <-unsubscribe:
				return
			}

		}
	}(stream)
	return unsubscribe
}

func Unsubscribe(c chan bool) {
	c <- true
}

func Post(a ...interface{}) {
	PostText(a...)
}

func PostText(a ...interface{}) {
	postMessage(Normal, a...)
}

func PostInfo(a ...interface{}) {
	postMessage(Info, a...)
}

func PostTitle(a ...interface{}) {
	postMessage(Title, a...)
}

func PostWarning(a ...interface{}) {
	postMessage(Warning, a...)
}

func PostError(a ...interface{}) {
	postMessage(Error, a...)
}

func postMessage(msgType MessageType, a ...interface{}) {
	message := Message{msgType, fmt.Sprint(a...), time.Now()}
	//fmt.Println("Reporting message:", message)
	msgProperty.Update(message)
}
