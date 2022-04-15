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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_append_tcr_event_to_csv_writer(t *testing.T) {
	testFlags := []struct {
		desc     string
		position int
		event    TcrEvent
		expected string
	}{
		{
			"timestamp in UTC",
			0,
			*aTcrEvent(withTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.UTC))),
			"2022-04-11 15:52:03",
		},
		{
			"timestamp not in UTC",
			0,
			*aTcrEvent(withTimestamp(time.Date(
				2022, 4, 11, 15, 52, 3, 0,
				time.FixedZone("UTC-7", -7*60*60)))),
			"2022-04-11 22:52:03",
		},
		{
			"modified source lines",
			1,
			*aTcrEvent(withModifiedSrcLines(2)),
			"2",
		},
		{
			"modified test lines",
			2,
			*aTcrEvent(withModifiedTestLines(25)),
			"25",
		},
		{
			"added test cases",
			3,
			*aTcrEvent(withAddedTestCases(3)),
			"3",
		},
		{
			"build passing",
			4,
			*aTcrEvent(withPassingBuild()),
			"true",
		},
		{
			"build failing",
			4,
			*aTcrEvent(withFailingBuild()),
			"false",
		},
		{
			"tests passing",
			5,
			*aTcrEvent(withPassingTests()),
			"true",
		},
		{
			"tests failing",
			5,
			*aTcrEvent(withFailingTests()),
			"false",
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			var b bytes.Buffer
			_ = appendEvent(tt.event, &b)
			str := strings.TrimSuffix(b.String(), "\n")
			fields := strings.Split(str, ",")
			assert.Equal(t, tt.expected, fields[tt.position])
		})
	}
}

func Test_it_should_create_the_file_when_it_doesnt_exist(t *testing.T) {
	event := aTcrEvent(withTimestamp(time.Date(
		2022, 4, 11, 15, 52, 3, 0,
		time.UTC)))

	mapFs := afero.NewMemMapFs()
	dirError := mapFs.Mkdir("test", os.ModeDir)
	assert.Nil(t, dirError)

	metricsFile, fileError := mapFs.Create("test/metrics.csv")
	assert.Nil(t, fileError)

	err := appendEvent(*event, metricsFile)
	assert.Nil(t, err)

	file, dirError := mapFs.Open("test/metrics.csv")
	assert.Nil(t, dirError)

	b := make([]byte, 50)
	readCount, readError := file.Read(b)

	//fmt.Println(string(b))
	assert.Nil(t, readError)
	assert.NotEqualf(t, 0, readCount, "Empty file")
}
