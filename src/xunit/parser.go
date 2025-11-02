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

package xunit

import (
	"time"

	"github.com/mengdaming/go-junit"
)

// TestStats is the structure containing test Stats extracted from xUnit files
type TestStats struct {
	Total    int
	Passed   int
	Failed   int
	Skipped  int
	InError  int
	Run      int
	Duration time.Duration
}

// Parser encapsulates XUnit files parsing
type Parser struct {
	Stats *TestStats
}

// NewParser returns a new XUnit parser instance
func NewParser() *Parser {
	return &Parser{Stats: &TestStats{}}
}

func (p *Parser) getTotalTests() int {
	return p.Stats.Total
}

func (p *Parser) getTotalTestsPassed() int {
	return p.Stats.Passed
}

func (p *Parser) getTotalTestsFailed() int {
	return p.Stats.Failed
}

func (p *Parser) getTotalTestsSkipped() int {
	return p.Stats.Skipped
}

func (p *Parser) getTotalTestsInError() int {
	return p.Stats.InError
}

func (p *Parser) getTotalTestsRun() int {
	return p.Stats.Run
}

func (p *Parser) getTotalTestDuration() time.Duration {
	return p.Stats.Duration
}

func (p *Parser) parse(xunitData []byte) error {
	suites, err := ingest(xunitData)
	if err != nil {
		return err
	}
	p.extractData(suites)
	return nil
}

func (p *Parser) resetCounters() {
	p.Stats = &TestStats{}
}

// ParseDir parses all xUnit files in the provided directory
func (p *Parser) ParseDir(dir string) error {
	suites, err := ingestDir(dir)
	if err != nil {
		return err
	}
	p.extractData(suites)
	return nil
}

func (p *Parser) extractData(suites []junit.Suite) {
	p.resetCounters()
	for _, suite := range suites {
		p.Stats.Total += suite.Totals.Tests
		p.Stats.Passed += suite.Totals.Passed
		p.Stats.Failed += suite.Totals.Failed
		p.Stats.Skipped += suite.Totals.Skipped
		p.Stats.InError += suite.Totals.Error
		p.Stats.Duration += suite.Totals.Duration
		p.Stats.Run += suite.Totals.Passed + suite.Totals.Failed + suite.Totals.Error
	}
}
