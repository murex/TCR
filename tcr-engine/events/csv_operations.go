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
	"encoding/csv"
	"fmt"
	"github.com/murex/tcr/tcr-engine/config"
	"github.com/murex/tcr/tcr-engine/report"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const timeLayoutFormat = "2006-01-02 15:04:05"

const eventLogFileName = "event-log.csv"

// AppendEventToLogFile appends a TCR event to the TCR event log file
func AppendEventToLogFile(e TcrEvent) {
	eventLogFilePath := filepath.Join(config.GetConfigDirPath(), eventLogFileName)
	//report.PostInfo("Metrics file: ", eventLogFilePath)
	file, err := os.OpenFile(eventLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		report.PostWarning(err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			report.PostWarning(err)
		}
	}(file)

	eventErr := appendEvent(e, file)
	if eventErr != nil {
		report.PostWarning(err)
		return
	}
}

// AppendEventToLogFile2 appends a TCR event to the TCR event log file
func AppendEventToLogFile2(e TcrEvent) {
	eventLogFilePath := filepath.Join("test", "test.csv")
	fmt.Println("Metrics file: ", eventLogFilePath)
	file, err := AppFs.OpenFile(eventLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		fmt.Println(err)
		report.PostWarning(err)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
			report.PostWarning(err)
			return
		}
	}()

	eventErr := appendEvent(e, file)
	if eventErr != nil {
		fmt.Println(err)
		report.PostWarning(err)
		return
	}
}

func appendEvent(event TcrEvent, out io.Writer) (err error) {
	w := csv.NewWriter(out)
	if err = w.Write(toCsvRecord(event)); err != nil {
		return
	}
	// Write any buffered data to the underlying writer.
	w.Flush()
	err = w.Error()
	return
}

func toCsvRecord(event TcrEvent) []string {
	return []string{
		event.Timestamp.In(time.UTC).Format(timeLayoutFormat),
		strconv.Itoa(event.ModifiedSrcLines),
		strconv.Itoa(event.ModifiedTestLines),
		strconv.Itoa(event.AddedTestCases),
		strconv.FormatBool(event.BuildPassed),
		strconv.FormatBool(event.TestsPassed),
	}
}
