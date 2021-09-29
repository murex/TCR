package report

import (
	"fmt"
	"github.com/imkira/go-observer"
	"time"
)

// MessageType type used for message characterization
type MessageType int

// List of possible values for MessageType field
const (
	Normal MessageType = iota
	Info
	Title
	Warning
	Error
	Event
)

// Message is the placeholder for any reported message
type Message struct {
	Type      MessageType
	Text      string
	Timestamp time.Time
}

var (
	msgProperty = observer.NewProperty(Message{Type: Normal, Text: ""})
)

// Subscribe allows a listener to subscribe to any posted message through the reporter.
// onReport() will be called every time a new message is posted. The returned channel
// shall be kept by the listener as this channel will be used for unsubscription
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

// Unsubscribe unsubscribes the listener associated to the provided channel from being notified
// of any further reported messages
func Unsubscribe(c chan bool) {
	c <- true
}

// Post posts some text for reporting. This is actually the same as PostText()
func Post(a ...interface{}) {
	PostText(a...)
}

// PostText posts some text for reporting
func PostText(a ...interface{}) {
	postMessage(Normal, a...)
}

// PostInfo posts an information message for reporting
func PostInfo(a ...interface{}) {
	postMessage(Info, a...)
}

// PostTitle posts a title message for reporting
func PostTitle(a ...interface{}) {
	postMessage(Title, a...)
}

// PostWarning posts a warning message for reporting
func PostWarning(a ...interface{}) {
	postMessage(Warning, a...)
}

// PostError posts an error message for reporting
func PostError(a ...interface{}) {
	postMessage(Error, a...)
}

// PostEvent posts an event message for reporting
func PostEvent(a ...interface{}) {
	postMessage(Event, a...)
}

func postMessage(msgType MessageType, a ...interface{}) {
	message := Message{msgType, fmt.Sprint(a...), time.Now()}
	//fmt.Println("Reporting message:", message)
	msgProperty.Update(message)
}
