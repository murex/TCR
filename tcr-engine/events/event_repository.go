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

package events

import "github.com/spf13/afero"

// AppFs is the singleton referring to the filesystem being used
var AppFs afero.Fs

// TcrEventRepository is the interface for a repository where we store TCR events
type TcrEventRepository interface {
	Add(event TcrEvent)
	GetLast() TcrEvent
}

// EventRepository is the TcrEventRepository instance
var EventRepository TcrEventRepository

func init() {
	EventRepository = &TcrEventFileRepository{}
	AppFs = afero.NewOsFs()
}

// TcrEventInMemoryRepository is an in-memory implementation of a TCR Event repository
type TcrEventInMemoryRepository struct {
	event TcrEvent
}

// GetLast returns the last event stored in TCR Event repository
func (t *TcrEventInMemoryRepository) GetLast() TcrEvent {
	return t.event
}

// Add stores a TCR event in TCR Event repository
func (t *TcrEventInMemoryRepository) Add(event TcrEvent) {
	t.event = event
}

// TcrEventFileRepository is a filesystem implementation of a TCR Event repository
type TcrEventFileRepository struct {
	filename string
	event    TcrEvent
}

// NewTcrEventFileRepository creates a new TCR Event file repository
func NewTcrEventFileRepository(filename string) TcrEventRepository {
	return &TcrEventFileRepository{filename: filename}
}

// GetLast returns the last event stored in TCR Event repository
func (t *TcrEventFileRepository) GetLast() TcrEvent {
	events := ReadEventLogFile()
	return events[len(events)-1]
}

// Add stores a TCR event in TCR Event repository
func (t *TcrEventFileRepository) Add(event TcrEvent) {
	t.event = event
	AppendEventToLogFile(t.event)
}
