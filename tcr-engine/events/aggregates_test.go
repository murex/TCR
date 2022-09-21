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
	"testing"
	"time"
)

func Test_aggregates_getters(t *testing.T) {
	testFlags := []struct {
		desc        string
		aggregates  Aggregates
		expectedMin interface{}
		expectedAvg interface{}
		expectedMax interface{}
	}{
		{
			"int aggregates",
			IntAggregates{1, 2.5, 4},
			1,
			2.5,
			4,
		},
		{
			"duration aggregates",
			DurationAggregates{
				500 * time.Millisecond,
				1*time.Second + 500*time.Millisecond,
				5 * time.Second,
			},
			500 * time.Millisecond,
			1*time.Second + 500*time.Millisecond,
			5 * time.Second,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedMin, tt.aggregates.Min())
			assert.Equal(t, tt.expectedAvg, tt.aggregates.Avg())
			assert.Equal(t, tt.expectedMax, tt.aggregates.Max())
		})
	}
}
