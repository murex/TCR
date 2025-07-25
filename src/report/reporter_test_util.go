//go:build test_helper

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

package report

import (
	"sync"
	"time"

	"github.com/murex/tcr/report/role_event"
	"github.com/murex/tcr/report/text"
	"github.com/murex/tcr/report/timer_event"
)

type messageFilter func(msg Message) bool

// Sniffer is a test utility allowing to track captured sent through TCR reporter
type Sniffer struct {
	reportingChannel chan bool
	filters          []messageFilter
	captured         []Message
	reporter         *Reporter
	mutex            sync.RWMutex
}

// NewSniffer creates a new instance of report sniffer, with filtering.
// If more than one filter is provided, the sniffer keeps all messages satisfying at least
// one of the filters
func NewSniffer(filters ...messageFilter) *Sniffer {
	return NewSnifferWithReporter(nil, filters...)
}

// NewSnifferWithReporter creates a new instance of report sniffer for a specific reporter instance.
// If reporter is nil, it uses the default reporter.
// If more than one filter is provided, the sniffer keeps all messages satisfying at least
// one of the filters
func NewSnifferWithReporter(reporter *Reporter, filters ...messageFilter) *Sniffer {
	sniffer := Sniffer{reporter: reporter}
	sniffer.addFilters(filters...)
	sniffer.Start()
	return &sniffer
}

func (sniffer *Sniffer) addFilters(filters ...messageFilter) {
	for _, filter := range filters {
		sniffer.addFilter(filter)
	}
}

func (sniffer *Sniffer) addFilter(filter messageFilter) {
	sniffer.filters = append(sniffer.filters, filter)
}

// Start tells the sniffer to start
func (sniffer *Sniffer) Start() {
	if sniffer.reporter != nil {
		sniffer.reportingChannel = sniffer.reporter.Subscribe(sniffer)
	} else {
		sniffer.reportingChannel = Subscribe(sniffer)
	}
}

func (sniffer *Sniffer) sniff(msg Message) {
	sniffer.mutex.Lock()
	defer sniffer.mutex.Unlock()

	if len(sniffer.filters) == 0 {
		// If no filter is set, we keep all messages
		sniffer.captured = append(sniffer.captured, msg)
		return
	}
	// Otherwise we keep any message satisfying at least one filter
	for _, filter := range sniffer.filters {
		if filter(msg) {
			sniffer.captured = append(sniffer.captured, msg)
			break
		}
	}
}

// ReportSimple reports simple messages
func (sniffer *Sniffer) ReportSimple(emphasis bool, payload text.Message) {
	sniffer.sniff(NewMessage(MessageType{Normal, emphasis}, payload))
}

// ReportInfo reports info messages
func (sniffer *Sniffer) ReportInfo(emphasis bool, payload text.Message) {
	sniffer.sniff(NewMessage(MessageType{Info, emphasis}, payload))
}

// ReportTitle reports title messages
func (sniffer *Sniffer) ReportTitle(emphasis bool, payload text.Message) {
	sniffer.sniff(NewMessage(MessageType{Title, emphasis}, payload))
}

// ReportSuccess reports warning messages
func (sniffer *Sniffer) ReportSuccess(emphasis bool, payload text.Message) {
	sniffer.sniff(NewMessage(MessageType{Success, emphasis}, payload))
}

// ReportWarning reports warning messages
func (sniffer *Sniffer) ReportWarning(emphasis bool, payload text.Message) {
	sniffer.sniff(NewMessage(MessageType{Warning, emphasis}, payload))
}

// ReportError reports error messages
func (sniffer *Sniffer) ReportError(emphasis bool, payload text.Message) {
	sniffer.sniff(NewMessage(MessageType{Error, emphasis}, payload))
}

// ReportRoleEvent reports role event messages
func (sniffer *Sniffer) ReportRoleEvent(emphasis bool, payload role_event.Message) {
	sniffer.sniff(NewMessage(MessageType{RoleEvent, emphasis}, payload))
}

// ReportTimerEvent reports timer event messages
func (sniffer *Sniffer) ReportTimerEvent(emphasis bool, payload timer_event.Message) {
	sniffer.sniff(NewMessage(MessageType{TimerEvent, emphasis}, payload))
}

// Stop tells the sniffer to stop
func (sniffer *Sniffer) Stop() {
	if sniffer.reportingChannel != nil {
		// Give a small moment for any pending messages to be processed
		time.Sleep(10 * time.Millisecond)

		// Use select to avoid blocking if channel is full
		select {
		case sniffer.reportingChannel <- true:
		case <-time.After(10 * time.Millisecond):
			// Channel might be blocked, that's ok
		}
	}
}

// FIXME SMELL: Duplicated assertion code in tests. Consider improving the reporter fake as a real mock with
//  check method and better error reporting, ex: report.assertWarning... instead of setting up the sniffer
//  then asserting count matches

// GetAllMatches returns a slice containing all matching messages captured by the sniffer
func (sniffer *Sniffer) GetAllMatches() []Message {
	sniffer.mutex.RLock()
	defer sniffer.mutex.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]Message, len(sniffer.captured))
	copy(result, sniffer.captured)
	return result
}

// GetMatchCount returns the number of matching messages captured by the sniffer
func (sniffer *Sniffer) GetMatchCount() int {
	sniffer.mutex.RLock()
	defer sniffer.mutex.RUnlock()
	return len(sniffer.captured)
}

// Test utility functions for isolated reporter testing

// TestWithIsolatedReporter runs a test function with an isolated reporter instance.
// This ensures that messages from the test don't interfere with other tests.
// The function restores the original default reporter when done.
func TestWithIsolatedReporter(testFunc func(reporter *Reporter, sniffer *Sniffer)) {
	// Create an isolated reporter for this test
	isolatedReporter := NewReporter()

	// Save the original default reporter
	originalReporter := SetDefaultReporter(isolatedReporter)

	// Create a sniffer for the isolated reporter
	sniffer := NewSnifferWithReporter(isolatedReporter)

	// Ensure cleanup happens even if test panics
	defer func() {
		sniffer.Stop()
		SetDefaultReporter(originalReporter)
	}()

	// Run the test function
	testFunc(isolatedReporter, sniffer)

	// Give a small moment to ensure all messages are processed
	time.Sleep(5 * time.Millisecond)
}

// TestWithIsolatedReporterAndFilters runs a test function with an isolated reporter instance and message filters.
// This ensures that messages from the test don't interfere with other tests.
// The function restores the original default reporter when done.
func TestWithIsolatedReporterAndFilters(testFunc func(reporter *Reporter, sniffer *Sniffer), filters ...messageFilter) {
	// Create an isolated reporter for this test
	isolatedReporter := NewReporter()

	// Save the original default reporter
	originalReporter := SetDefaultReporter(isolatedReporter)

	// Create a sniffer for the isolated reporter with filters
	sniffer := NewSnifferWithReporter(isolatedReporter, filters...)

	// Ensure cleanup happens even if test panics
	defer func() {
		sniffer.Stop()
		SetDefaultReporter(originalReporter)
	}()

	// Run the test function
	testFunc(isolatedReporter, sniffer)
}
