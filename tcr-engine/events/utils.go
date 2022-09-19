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

// MinMaxAvg is a convenience type for returning min, max and average values altogether
type MinMaxAvg interface {
	Min() interface{}
	Avg() interface{}
	Max() interface{}
}

// MinMaxAvgDuration is the MinMaxAvg implementation for time.Duration type
type MinMaxAvgDuration struct {
	min time.Duration
	avg time.Duration
	max time.Duration
}

// Min returns the min value
func (m MinMaxAvgDuration) Min() interface{} {
	return m.min
}

// Avg returns the average value
func (m MinMaxAvgDuration) Avg() interface{} {
	return m.avg
}

// Max returns the max value
func (m MinMaxAvgDuration) Max() interface{} {
	return m.max
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
