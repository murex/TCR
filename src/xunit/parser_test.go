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
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_retrieve_xunit_test_counters(t *testing.T) {
	testFlags := []struct {
		desc              string
		extractorFunction func(parser *Parser) int
		expected          int
	}{
		{
			"total tests",
			func(p *Parser) int { return p.getTotalTests() },
			sampleTotalsSuite0.Tests + sampleTotalsSuite1.Tests,
		},
		{
			"total tests passed",
			func(p *Parser) int { return p.getTotalTestsPassed() },
			sampleTotalsSuite0.Passed + sampleTotalsSuite1.Passed,
		},
		{
			"total tests failed",
			func(p *Parser) int { return p.getTotalTestsFailed() },
			sampleTotalsSuite0.Failed + sampleTotalsSuite1.Failed,
		},
		{
			"total tests in error",
			func(p *Parser) int { return p.getTotalTestsInError() },
			sampleTotalsSuite0.Error + sampleTotalsSuite1.Error,
		},
		{
			"total tests skipped",
			func(p *Parser) int { return p.getTotalTestsSkipped() },
			sampleTotalsSuite0.Skipped + sampleTotalsSuite1.Skipped,
		},
		{
			"total tests run",
			func(p *Parser) int { return p.getTotalTestsRun() },
			sampleTotalsSuite0.Passed + sampleTotalsSuite1.Passed +
				sampleTotalsSuite0.Failed + sampleTotalsSuite1.Failed +
				sampleTotalsSuite0.Error + sampleTotalsSuite1.Error,
		},
	}

	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			parser := NewParser()
			_ = parser.parse(xunitSample)
			assert.Equal(t, tt.expected, tt.extractorFunction(parser))
		})
	}
}

func Test_retrieve_xunit_test_duration(t *testing.T) {
	parser := NewParser()
	_ = parser.parse(xunitSample)
	assert.Equal(t, sampleTotalsSuite0.Duration+sampleTotalsSuite1.Duration, parser.getTotalTestDuration())
}

func Test_parsing_invalid_data(t *testing.T) {
	testFlags := []struct {
		desc        string
		xunitData   []byte
		expectError bool
	}{
		{"no data",
			nil,
			false,
		},
		{"empty data",
			[]byte(""),
			false,
		},
		{"header only data",
			[]byte(`<?xml version="1.0" encoding="UTF-8"?>`),
			false,
		},
		{"truncated data",
			[]byte(`<?xml version="1.0" encoding="UTF-8"?><testsuites`),
			true,
		},
	}
	for _, tt := range testFlags {
		t.Run(tt.desc, func(t *testing.T) {
			err := NewParser().parse(tt.xunitData)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_parse_dir(t *testing.T) {
	appFs = afero.NewMemMapFs()
	_ = appFs.Mkdir("build", os.ModeDir)
	_ = afero.WriteFile(appFs, "build/sample1.xml", xunitSample, 0644)
	_ = afero.WriteFile(appFs, "build/sample2.xml", xunitSample, 0644)
	_ = afero.WriteFile(appFs, "build/sample3.xml", xunitSample, 0644)

	parser := NewParser()
	assert.NoError(t, parser.ParseDir("build"))
	assert.Equal(t, 3*(sampleTotalsSuite0.Tests+sampleTotalsSuite1.Tests), parser.getTotalTests())
	assert.Equal(t, 3*(sampleTotalsSuite0.Duration+sampleTotalsSuite1.Duration), parser.getTotalTestDuration())
}

func Test_parse_dir_with_error(t *testing.T) {
	appFs = afero.NewMemMapFs()
	_ = appFs.Mkdir("build", os.ModeDir)
	_ = afero.WriteFile(appFs, "build/sample1.xml", xunitSample, 0644)
	_ = afero.WriteFile(appFs, "build/sample2.xml",
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><testsuites`), 0644)

	parser := NewParser()
	assert.Error(t, parser.ParseDir("build"))
	assert.Equal(t, 0, parser.getTotalTests())
	assert.Equal(t, 0, parser.getTotalTestsPassed())
	assert.Equal(t, 0, parser.getTotalTestsFailed())
	assert.Equal(t, 0, parser.getTotalTestsInError())
	assert.Equal(t, 0, parser.getTotalTestsSkipped())
	assert.Equal(t, 0, parser.getTotalTestsRun())
	assert.Equal(t, time.Duration(0), parser.getTotalTestDuration())
}
