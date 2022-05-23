package charts

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/murex/tcr/tcr-engine/events"
	"github.com/pkg/browser"
	"path/filepath"
	"sort"
	"time"
)

// Data is a data record
type Data struct {
	Key   time.Time
	Value int
}

const htmlFile = "tcr-metrics.html"

// DataList contains the list of all records
type DataList []Data

// Len returns the size of DataList
func (d DataList) Len() int {
	return len(d)
}

// Swap swaps the records at indices i and j
func (d DataList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less returns the subtraction of values at indices i and j
func (d DataList) Less(i, j int) bool {
	return d[i].Value < d[j].Value
}

// Browse loads the CSV file contents, generates an HTML page, and browses it
func Browse() {
	e := events.ReadEventLogFile()
	sortedData := mapData(formatRecords(e))
	createChart(sortedData)

	err := browser.OpenFile(filepath.Join("", htmlFile))
	if err != nil {
		fmt.Println(err)
	}
}

func formatRecords(records []events.TcrEvent) (timestamp []time.Time, testsRun []int) {
	for _, e := range records {
		timestamp = append(timestamp, e.Timestamp)
		testsRun = append(testsRun, e.TotalTestsRun)
	}
	return
}

func mapData(timeStamps []time.Time, testRun []int) DataList {
	dataMap := map[time.Time]int{}
	for index, value := range timeStamps {
		dataMap[value] = testRun[index]
	}
	data := make(DataList, len(dataMap))
	iterator := 0
	for k, v := range dataMap {
		data[iterator] = Data{k, v}
		iterator++
	}
	sort.Sort(data)
	return data
}

func createChart(sortedData DataList) {
	chart := charts.NewLine()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "TCR log stats",
		Subtitle: "Total tests run per save",
	}))
	chart.SetXAxis(generateXAxisItems(sortedData)).AddSeries("Tests Run", generateLineItems(sortedData))
	f, _ := events.AppFs.Create(htmlFile)
	err := chart.Render(f)
	if err != nil {
		fmt.Println(err)
	}
}

func formatTimestamp(t time.Time) string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func generateXAxisItems(data DataList) (values []string) {
	for _, d := range data {
		values = append(values, formatTimestamp(d.Key))
	}
	return
}

func generateBarItems(data DataList) (items []opts.BarData) {
	for _, d := range data {
		items = append(items, opts.BarData{Value: d.Value})
	}
	return items
}

func generateLineItems(data DataList) (items []opts.LineData) {
	for _, d := range data {
		items = append(items, opts.LineData{Value: d.Value})
	}
	return items
}

//func generateKLineItems(data DataList) (items []opts.KlineData) {
//	for _, d := range data {
//		items = append(items, opts.KlineData{Value: d.Value})
//	}
//	return items
//}
