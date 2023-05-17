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
	"github.com/mengdaming/go-junit"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var xunitSample = []byte(`
    <?xml version="1.0" encoding="UTF-8"?>
    <testsuites>
        <testsuite name="JUnitXmlReporter" errors="0" tests="0" failures="0" time="0" timestamp="2013-05-24T10:23:58" />
        <testsuite name="JUnitXmlReporter.constructor" errors="0" skipped="1" tests="3" failures="1" time="0.006" timestamp="2013-05-24T10:23:58">
            <testcase classname="JUnitXmlReporter.constructor" name="should default path to an empty string" time="0.006">
                <failure message="test failure">Assertion failed</failure>
            </testcase>
            <testcase classname="JUnitXmlReporter.constructor" name="should default consolidate to true" time="0">
                <skipped />
            </testcase>
            <testcase classname="JUnitXmlReporter.constructor" name="should default useDotNotation to true" time="0" />
        </testsuite>
    </testsuites>
`)

var (
	sampleTotalsSuite0 = junit.Totals{Tests: 0, Passed: 0, Skipped: 0, Failed: 0, Error: 0, Duration: 0}
	sampleTestStatus0  []junit.Status
	sampleTotalsSuite1 = junit.Totals{Tests: 3, Passed: 1, Skipped: 1, Failed: 1, Error: 0, Duration: 6 * time.Millisecond}
	sampleTestStatus1  = []junit.Status{junit.StatusFailed, junit.StatusSkipped, junit.StatusPassed}
)

func assertXunitSampleData(t *testing.T, suites []junit.Suite, nbSamples int) {
	t.Helper()
	assert.Equal(t, 2*nbSamples, len(suites))

	for i := 0; i < nbSamples; i++ {
		suite0 := suites[2*i+0]
		assert.Equal(t, sampleTotalsSuite0, suite0.Totals)
		assert.Equal(t, len(sampleTestStatus0), len(suite0.Tests))

		suite1 := suites[2*i+1]
		assert.Equal(t, sampleTotalsSuite1, suite1.Totals)
		assert.Equal(t, len(sampleTestStatus1), len(suite1.Tests))
		assert.Equal(t, sampleTestStatus1[0], suite1.Tests[0].Status)
		assert.Equal(t, sampleTestStatus1[1], suite1.Tests[1].Status)
		assert.Equal(t, sampleTestStatus1[2], suite1.Tests[2].Status)
	}
}

func Test_ingest_wrapper(t *testing.T) {
	suites, err := ingest(xunitSample)
	if assert.NoError(t, err) {
		assertXunitSampleData(t, suites, 1)
	}
}

func Test_ingest_file_wrapper(t *testing.T) {
	appFs = afero.NewMemMapFs()
	_ = appFs.Mkdir("build", os.ModeDir)
	_ = afero.WriteFile(appFs, "build/sample.xml", xunitSample, 0644)
	suites, err := ingestFile("build/sample.xml")
	if assert.NoError(t, err) {
		assertXunitSampleData(t, suites, 1)
	}
}

func Test_ingest_file_wrapper_on_error(t *testing.T) {
	suites, err := ingestFile("missing-file.xml")
	assert.Error(t, err)
	assert.Zero(t, suites)
}

func Test_ingest_files_wrapper(t *testing.T) {
	appFs = afero.NewMemMapFs()
	_ = appFs.Mkdir("build", os.ModeDir)
	_ = afero.WriteFile(appFs, "build/sample1.xml", xunitSample, 0644)
	_ = afero.WriteFile(appFs, "build/sample2.xml", xunitSample, 0644)
	suites, err := ingestFiles([]string{"build/sample1.xml", "build/sample2.xml"})
	if assert.NoError(t, err) {
		assertXunitSampleData(t, suites, 2)
	}
}

func Test_ingest_files_wrapper_on_error(t *testing.T) {
	appFs = afero.NewMemMapFs()
	_ = appFs.Mkdir("build", os.ModeDir)
	_ = afero.WriteFile(appFs, "build/sample.xml", xunitSample, 0644)
	suites, err := ingestFiles([]string{"build/sample.xml", "missing-file.xml"})
	assert.Error(t, err)
	assert.Zero(t, suites)
}

func Test_ingest_dir_wrapper(t *testing.T) {
	appFs = afero.NewMemMapFs()
	_ = appFs.Mkdir("build", os.ModeDir)
	_ = afero.WriteFile(appFs, "build/sample1.xml", xunitSample, 0644)
	_ = afero.WriteFile(appFs, "build/sample2.xml", xunitSample, 0644)
	_ = afero.WriteFile(appFs, "build/sample3.xml", xunitSample, 0644)
	suites, err := ingestDir("build")
	if assert.NoError(t, err) {
		assertXunitSampleData(t, suites, 3)
	}
}

func Test_ingest_dir_wrapper_on_error(t *testing.T) {
	appFs = afero.NewMemMapFs()
	suites, err := ingestDir("build")
	assert.Error(t, err)
	assert.Zero(t, suites)
}

func Test_ingest_dir_wrapper_with_symbolic_link(t *testing.T) {
	appFs = afero.NewOsFs()
	tempDir, errTempDir := afero.TempDir(appFs, "", "tcr-xunit-test")
	if errTempDir != nil {
		t.Fatal(errTempDir)
	}
	t.Cleanup(func() {
		_ = appFs.RemoveAll(tempDir)
	})
	realDir := filepath.Join(tempDir, "build")
	linkDir := filepath.Join(tempDir, "link")
	_ = appFs.Mkdir(realDir, 0755)
	_ = afero.WriteFile(appFs, filepath.Join(realDir, "sample1.xml"), xunitSample, 0644)
	errSymLink := appFs.(afero.Linker).SymlinkIfPossible(realDir, linkDir)
	if errSymLink != nil {
		// Note: this test is skipped on windows due to its crappy handling of symlinks
		// (requires elevated privileges to create a symbolic link)
		t.Skip("symbolic links not supported: ", errSymLink)
	}
	suites, err := ingestDir(linkDir)
	if assert.NoError(t, err) {
		assertXunitSampleData(t, suites, 1)
	}
}
