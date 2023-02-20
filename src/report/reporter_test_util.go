//go:build test_helper

/*
Copyright (c) 2022 Murex

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

import "time"

type messageFilter func(msg Message) bool

// Sniffer is a test utility allowing to track captured sent through TCR reporter
type Sniffer struct {
	reportingChannel chan bool
	filters          []messageFilter
	captured         []Message
}

// NewSniffer creates a new instance of report sniffer, with filtering.
// If more than one filter is provided, the sniffer keeps all messages satisfying at least
// one of the filters
func NewSniffer(filters ...messageFilter) *Sniffer {
	sniffer := Sniffer{}
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
	sniffer.reportingChannel = Subscribe(sniffer)
}

func (sniffer *Sniffer) sniff(msg Message) {
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
func (sniffer *Sniffer) ReportSimple(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Normal, emphasis}, a...))
}

// ReportInfo reports info messages
func (sniffer *Sniffer) ReportInfo(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Info, emphasis}, a...))
}

// ReportTitle reports title messages
func (sniffer *Sniffer) ReportTitle(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Title, emphasis}, a...))
}

// ReportTimer reports timer messages
func (sniffer *Sniffer) ReportTimer(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Timer, emphasis}, a...))
}

// ReportSuccess reports warning messages
func (sniffer *Sniffer) ReportSuccess(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Success, emphasis}, a...))
}

// ReportWarning reports warning messages
func (sniffer *Sniffer) ReportWarning(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Warning, emphasis}, a...))
}

// ReportError reports error messages
func (sniffer *Sniffer) ReportError(emphasis bool, a ...interface{}) {
	sniffer.sniff(NewMessage(MessageType{Error, emphasis}, a...))
}

// Stop tells the sniffer to stop
func (sniffer *Sniffer) Stop() {
	// Micro-timer to give a chance to reporters to process their messages
	// before the sniffer stops listening
	time.Sleep(1 * time.Millisecond)
	if sniffer.reportingChannel != nil {
		Unsubscribe(sniffer.reportingChannel)
	}
}

// GetAllMatches returns a slice containing all matching messages captured by the sniffer
func (sniffer *Sniffer) GetAllMatches() []Message {
	return sniffer.captured
}

// GetMatchCount returns the number of matching messages captured by the sniffer
func (sniffer *Sniffer) GetMatchCount() int {
	return len(sniffer.captured)
}
