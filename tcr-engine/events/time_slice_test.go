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

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func Test_time_slice_sorting(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(1 * time.Second)
	earlier := now.Add(-1 * time.Second)
	testFlags := []struct {
		desc     string
		input    TimeSlice
		expected TimeSlice
	}{
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty slice",
			TimeSlice{},
			TimeSlice{},
		},
		{
			"1 timestamp",
			TimeSlice{now},
			TimeSlice{now},
		},
		{
			"2 ordered timestamps",
			TimeSlice{now, later},
			TimeSlice{now, later},
		},
		{
			"2 inverted timestamps",
			TimeSlice{later, now},
			TimeSlice{now, later},
		},
		{
			"3 timestamps",
			TimeSlice{later, earlier, now},
			TimeSlice{earlier, now, later},
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			sort.Sort(tt.input)
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}
