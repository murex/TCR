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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_value_and_ratio_getters(t *testing.T) {
	testFlags := []struct {
		desc               string
		valueAndRatio      ValueAndRatio
		expectedValue      interface{}
		expectedPercentage int
	}{
		{
			"int value and ratio",
			IntValueAndRatio{10, 34},
			10,
			34,
		},
		{
			"duration value and ratio",
			DurationValueAndRatio{1 * time.Second, 28},
			1 * time.Second,
			28,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expectedValue, tt.valueAndRatio.Value(), "value")
			assert.Equal(t, tt.expectedPercentage, tt.valueAndRatio.Percentage(), "percentage")
		})
	}
}
