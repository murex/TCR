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

package metrics

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_write_one_timestamp_to_csv(t *testing.T) {
	event := aTcrEvent(
		withTimestamp(
			time.Date(2022, 4, 11, 15, 52, 3, 0, time.UTC)))
	var b bytes.Buffer
	appendEvent(event, &b)
	assert.Equal(t, "2022-04-11 15:52:03\n", b.String())
}

func Test_the_time_stamp_should_be_saved_as_UTC_time(t *testing.T) {
	event := aTcrEvent(
		withTimestamp(
			time.Date(2022, 4, 11, 15, 52, 3, 0,
				time.FixedZone("UTC-7", -7*60*60))))
	var b bytes.Buffer
	appendEvent(event, &b)
	assert.Equal(t, "2022-04-11 22:52:03\n", b.String())
}
