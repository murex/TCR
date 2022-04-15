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
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

const timeLayoutFormat = "2006-01-02 15:04:05"

//const metricsFileName = "./TCR_Metrics.csv"

//func appendEventToMetricsFile(event TcrEvent) {
//	file, err := os.OpenFile(metricsFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
//	if err != nil {
//		log.Fatalln("error opening the metrics file: ", err)
//	}
//	defer file.Close()
//
//	eventErr := appendEvent(event, file)
//	if eventErr != nil {
//		return
//	}
//}

func appendEvent(event TcrEvent, out io.Writer) (err error) {
	w := csv.NewWriter(out)
	if err = w.Write(
		[]string{
			event.timestamp.In(time.UTC).Format(timeLayoutFormat),
			strconv.Itoa(event.modifiedSrcLines),
			strconv.Itoa(event.modifiedTestLines),
			strconv.Itoa(event.addedTestCases),
			strconv.FormatBool(event.buildPassed),
			strconv.FormatBool(event.testsPassed),
		}); err != nil {
		return
	}
	// Write any buffered data to the underlying writer.
	w.Flush()
	err = w.Error()
	return
}
