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
	"github.com/murex/tcr/tcr-engine/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	config.DirPathGetter = func() string {
		return "test"
	}
}

func Test_add_and_get_event_to_in_memory_repository(t *testing.T) {
	repository := &TcrEventInMemoryRepository{}
	event := *ATcrEvent()
	repository.Add(event)
	assert.Equal(t, event, repository.Get())
}

func Test_add_a_single_event_to_file_repository(t *testing.T) {
	repository, filePath := setUpFileRepository()

	tcrEvent := ATcrEvent()
	repository.Add(*tcrEvent)

	eventLogBytes, _ := afero.ReadFile(AppFs, filePath)
	assert.Equal(t, "0,0,0,0,0,0,0\n", strings.Trim(string(eventLogBytes), " "))
}

func Test_gets_a_single_event_from_file_repository(t *testing.T) {
	repository, filePath := setUpFileRepository()

	eventStr := "0,0,0,0,0,0,0\n"

	_ = afero.WriteFile(AppFs, filePath, []byte(eventStr), os.ModePerm)

	tcrEvent := ATcrEvent()

	assert.Equal(t, *tcrEvent, repository.Get())
}

func setUpFileRepository() (TcrEventRepository, string) {
	AppFs = afero.NewMemMapFs()
	eventFilePath := getEventLogFileName()

	_ = AppFs.Mkdir(config.DirPathGetter(), os.ModeDir)
	_, _ = AppFs.Create(eventFilePath)

	return NewTcrEventFileRepository(eventFilePath), eventFilePath
}

func getEventLogFileName() string {
	filename := "event-log.csv"
	return filepath.Join(config.DirPathGetter(), filename)
}
