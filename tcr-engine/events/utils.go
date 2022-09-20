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

// ZeroTime is Unix starting time, e.g. Jan 1, 1970
var ZeroTime = time.Unix(0, 0).UTC()

// Aggregates is a convenience type for returning min, max, average, etc. altogether
type Aggregates interface {
	Min() interface{}
	Avg() interface{}
	Max() interface{}
}

// DurationAggregates is the Aggregates implementation for time.Duration type
type DurationAggregates struct {
	min time.Duration
	avg time.Duration
	max time.Duration
}

// Min returns the min duration value
func (da DurationAggregates) Min() interface{} {
	return da.min
}

// Avg returns the average duration value
func (da DurationAggregates) Avg() interface{} {
	return da.avg
}

// Max returns the max duration value
func (da DurationAggregates) Max() interface{} {
	return da.max
}

func inSeconds(duration time.Duration) int {
	return int(duration / time.Second)
}

func asPercentage(dividend, divisor int) int {
	return roundToClosestInt(100*dividend, divisor) //nolint:revive
}

func roundToClosestInt(dividend, divisor int) int {
	if divisor == 0 {
		return 0
	}
	return (dividend + (divisor / 2)) / divisor
}

// IntAggregates is the Aggregates implementation for int type
type IntAggregates struct {
	min int
	avg float32
	max int
}

// Min returns the min value
func (ia IntAggregates) Min() interface{} {
	return ia.min
}

// Avg returns the average value
func (ia IntAggregates) Avg() interface{} {
	return ia.avg
}

// Max returns the max value
func (ia IntAggregates) Max() interface{} {
	return ia.max
}

// ValueEvolution describes evolution of a value between first and last record
type ValueEvolution interface {
	From() interface{}
	To() interface{}
}

// IntValueEvolution implements ValueEvolution interface for an int value
type IntValueEvolution struct {
	from int
	to   int
}

// From returns the starting value of an IntValueEvolution instance
func (iv IntValueEvolution) From() interface{} {
	return iv.from
}

// To returns the ending value of an IntValueEvolution instance
func (iv IntValueEvolution) To() interface{} {
	return iv.to
}
