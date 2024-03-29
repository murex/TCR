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

// ValueEvolution describes evolution of a value between first and last record
type ValueEvolution interface {
	From() any
	To() any
}

// DurationValueEvolution implements ValueEvolution interface for a time.Duration value
type DurationValueEvolution struct {
	from time.Duration
	to   time.Duration
}

// From returns the starting value of a DurationValueEvolution instance
func (dve DurationValueEvolution) From() any {
	return dve.from
}

// To returns the ending value of a DurationValueEvolution instance
func (dve DurationValueEvolution) To() any {
	return dve.to
}

// IntValueEvolution implements ValueEvolution interface for an int value
type IntValueEvolution struct {
	from int
	to   int
}

// From returns the starting value of an IntValueEvolution instance
func (ive IntValueEvolution) From() any {
	return ive.from
}

// To returns the ending value of an IntValueEvolution instance
func (ive IntValueEvolution) To() any {
	return ive.to
}
