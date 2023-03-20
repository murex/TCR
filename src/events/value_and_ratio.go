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

import "time"

// ValueAndRatio provides a value with its associated percentage.
// The percentage is rounded to an int value (ex: percentage=50 means 50%)
type ValueAndRatio interface {
	Value() any
	Percentage() int
}

// DurationValueAndRatio implements ValueAndRatio interface for a time.Duration value
type DurationValueAndRatio struct {
	value      time.Duration
	percentage int
}

// Value returns the value of a DurationValueAndRatio instance
func (dvr DurationValueAndRatio) Value() any {
	return dvr.value
}

// Percentage returns the percentage of a DurationValueAndRatio instance
func (dvr DurationValueAndRatio) Percentage() int {
	return dvr.percentage
}

// IntValueAndRatio implements ValueAndRatio interface for an int value
type IntValueAndRatio struct {
	value      int
	percentage int
}

// Value returns the value of an IntValueAndRatio instance
func (ivr IntValueAndRatio) Value() any {
	return ivr.value
}

// Percentage returns the percentage of an IntValueAndRatio instance
func (ivr IntValueAndRatio) Percentage() int {
	return ivr.percentage
}
