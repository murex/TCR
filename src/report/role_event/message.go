/*
Copyright (c) 2024 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package role_event //nolint:revive

import (
	"fmt"
	"github.com/murex/tcr/role"
	"strings"
)

// Trigger represents what triggers a role event
type Trigger string

// List of possible Trigger values
const (
	TriggerStart Trigger = "start"
	TriggerEnd   Trigger = "end"
)

const separator = ":"

// Message contains a role event information
type Message struct {
	Trigger Trigger
	Role    role.Role
}

// WithEmphasis indicates whether the event message should be reported with emphasis flag
func (Message) WithEmphasis() bool {
	return false
}

// WrapMessage wraps a role event Message into a string
func WrapMessage(em Message) string {
	return fmt.Sprint(em.Role.Name(), separator, em.Trigger)
}

// UnwrapMessage unwraps a role event message string into a role event Message instance
func UnwrapMessage(message string) Message {
	parts := strings.Split(message, separator)
	return Message{
		Role:    role.FromName(parts[0]),
		Trigger: Trigger(parts[1]),
	}
}
