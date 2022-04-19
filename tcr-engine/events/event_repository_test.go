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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func Test_add_and_get_event_to_in_memory_repository(t *testing.T) {
	repository := &TcrEventInMemoryRepository{}
	event := *ATcrEvent()
	repository.Add(event)
	assert.Equal(t, event, repository.Get())
}

func Test_add_and_get_event_to_file_repository(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	filename := "test.csv"
	eventFilePath := filepath.Join("test", filename)

	_ = AppFs.Mkdir("test", os.ModeDir)
	_, _ = AppFs.Create(eventFilePath)

	repository := NewTcrEventFileRepository(eventFilePath)
	event := *ATcrEvent()
	repository.Add(event)

	file, _ := AppFs.Open(eventFilePath)
	b := make([]byte, 50)
	readCount, readError := file.Read(b)

	assert.Nil(t, readError)
	assert.NotEqualf(t, 0, readCount, "Empty file")

	assert.Equal(t, event, repository.Get())
}
