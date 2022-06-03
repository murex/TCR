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
	"github.com/joshdk/go-junit"
	"time"
)

// Parser encapsulates XUnit files parsing
type Parser struct {
	total    int
	passed   int
	failed   int
	skipped  int
	inError  int
	duration time.Duration
}

// NewParser returns a new XUnit parser instance
func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) parse(xunitData []byte) (err error) {
	p.total, p.passed, p.failed, p.skipped, p.inError, p.duration = 0, 0, 0, 0, 0, 0
	var suites []junit.Suite
	suites, err = Ingest(xunitData)
	if err != nil {
		return
	}
	for _, suite := range suites {
		p.total += suite.Totals.Tests
		p.passed += suite.Totals.Passed
		p.failed += suite.Totals.Failed
		p.skipped += suite.Totals.Skipped
		p.inError += suite.Totals.Error
		p.duration += suite.Totals.Duration
	}
	return
}

func (p *Parser) getTotalTests() int {
	return p.total
}

func (p *Parser) getTotalTestsPassed() int {
	return p.passed
}

func (p *Parser) getTotalTestsFailed() int {
	return p.failed
}

func (p *Parser) getTotalTestsSkipped() int {
	return p.skipped
}

func (p *Parser) getTotalTestsInError() int {
	return p.inError
}

func (p *Parser) getTotalTestsRun() int {
	return p.passed + p.failed + p.inError
}

func (p *Parser) getTotalTestDuration() time.Duration {
	return p.duration
}
