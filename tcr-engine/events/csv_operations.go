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
	"github.com/murex/tcr/tcr-engine/config"
	"github.com/murex/tcr/tcr-engine/report"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const timeLayoutFormat = "2006-01-02 15:04:05"

const eventLogFileName = "event-log.csv"

// AppendEventToLogFile appends a TCR event to the TCR event log file
func AppendEventToLogFile(e TcrEvent) {
	eventLogFile := openEventLogFile()

	if eventLogFile == nil {
		return
	}

	defer func() {
		err := eventLogFile.Close()
		if err != nil {
			report.PostWarning(err)
			return
		}
	}()

	eventErr := appendEvent(e, eventLogFile)
	if eventErr != nil {
		report.PostWarning(eventErr)
		return
	}
}

func openEventLogFile() afero.File {
	eventLogFilePath := filepath.Join(config.DirPathGetter(), eventLogFileName)
	file, err := AppFs.OpenFile(eventLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		report.PostWarning(err)
		return nil
	}
	return file
}

// ReadEventLogFile reads the content of the EventLog file
func ReadEventLogFile() TcrEvent {
	a := afero.Afero{
		Fs: AppFs,
	}
	eventLogBytes, err := a.ReadFile(filepath.Join(config.DirPathGetter(), eventLogFileName))

	if err != nil {
		report.PostWarning(err)
		return TcrEvent{}
	}

	fileContent := string(eventLogBytes)
	return toTcrEvent(fileContent)
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

func toTcrEvent(csvRecord string) TcrEvent {
	split := strings.Split(csvRecord, ",")
	parsedTime, _ := time.Parse(timeLayoutFormat, split[0])

	event := TcrEvent{
		Timestamp:         parsedTime,
		ModifiedSrcLines:  toInt(split[1]),
		ModifiedTestLines: toInt(split[2]),
		AddedTestCases:    toInt(split[3]),
		BuildPassed:       toBoolean(split[4]),
		TestsPassed:       toBoolean(split[5]),
	}
	return event
}

func toBoolean(value string) bool {
	spaceTrimmed := strings.Trim(value, "\n")
	trimmed := strings.Trim(spaceTrimmed, " ")
	parseBool, _ := strconv.ParseBool(trimmed)
	return parseBool
}

func toInt(value string) int {
	intValue, _ := strconv.Atoi(strings.Trim(value, " "))
	return intValue
}
