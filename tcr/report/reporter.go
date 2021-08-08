package report

import (
	"github.com/imkira/go-observer"
)

var (
	msgProperty = observer.NewProperty("")
)

func Register(onChange func(msg string)) {
	stream := msgProperty.Observe()

	val := stream.Value().(string)
	//fmt.Printf("initial value: %v\n", val)

	for {
		select {
		// wait for changes
		case <-stream.Changes():
			// advance to next value
			stream.Next()
			val = stream.Value().(string)
			//fmt.Printf("got new value: %v\n", val)
			onChange(val)
		}
	}
}

func Report(message string) {
	//fmt.Println("Reporting message:", message)
	msgProperty.Update(message)
}

