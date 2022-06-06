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
)

const timeLayoutFormat = "2006-01-02 15:04:05"

const eventLogFileName = "event-log.csv"

// AppendEventToLogFile appends a TCR event to the TCR event log file
func AppendEventToLogFile(e TcrEvent) {
	eventLogFile, err := openEventLogFile()

	if err != nil {
		report.PostWarning(err)
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

func openEventLogFile() (afero.File, error) {
	eventLogFilePath := filepath.Join(config.DirPathGetter(), eventLogFileName)
	return AppFs.OpenFile(eventLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
}

// ReadEventLogFile reads the content of the EventLog file
func ReadEventLogFile() TcrEvent {
	// TODO: Read the file line by line, then convert each line to a TCR Event
	eventLogBytes, err := afero.ReadFile(AppFs, filepath.Join(config.DirPathGetter(), eventLogFileName))

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
		strconv.Itoa(event.ModifiedSrcLines),
		strconv.Itoa(event.ModifiedTestLines),
		strconv.Itoa(event.TotalTestsRun),
		strconv.Itoa(event.TestsPassed),
		strconv.Itoa(event.TestsFailed),
		strconv.Itoa(event.TestsSkipped),
		strconv.Itoa(event.TestsWithErrors),
	}
}

func toTcrEvent(csvRecord string) TcrEvent {
	trimmedRecord := strings.Trim(csvRecord, "\n")
	recordFields := strings.Split(trimmedRecord, ",")

	event := TcrEvent{
		ModifiedSrcLines:  toInt(recordFields[0]),
		ModifiedTestLines: toInt(recordFields[1]),
		TotalTestsRun:     toInt(recordFields[2]),
		TestsPassed:       toInt(recordFields[3]),
		TestsFailed:       toInt(recordFields[4]),
		TestsSkipped:      toInt(recordFields[5]),
		TestsWithErrors:   toInt(recordFields[6]),
	}
	return event
}

func toTcrEventStatus(field string) TcrEventStatus {
	return TcrEventStatus(toInt(field))
}

func toInt(field string) int {
	intValue, _ := strconv.Atoi(strings.TrimSpace(field))
	return intValue
}
