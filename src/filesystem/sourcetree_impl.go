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

package filesystem

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/murex/tcr/report"
	"os"
	"path/filepath"
)

// SourceTreeImpl is the implementation of Source Tree interface
type SourceTreeImpl struct {
	baseDir string
	valid   bool
	watcher *fsnotify.Watcher
	matcher func(filename string) bool
}

// New creates a new instance of source tree implementation with a root directory set as dir.
// The method returns an error if the root directory does not exist or cannot be accessed.
func New(dir string) (SourceTree, error) {
	var impl = &SourceTreeImpl{}
	var err error
	impl.baseDir, err = checkDir(dir)
	if err != nil {
		return nil, err
	}
	impl.valid = true
	return impl, nil
}

// IsValid indicates that the source tree instance is valid
func (st *SourceTreeImpl) IsValid() bool {
	return st.valid
}

// checkDir returns the absolute path for the provided directory.
// Returns an error if the directory cannot be accessed or is not a directory
func checkDir(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err == nil {
		info, err := os.Stat(absPath)
		if err != nil {
			return "", errors.New("cannot access " + absPath)
		}
		if !info.IsDir() {
			return "", errors.New(absPath + " exists but is not a directory")
		}
	}
	return absPath, nil
}

// GetBaseDir returns the base directory for the source tree instance
func (st *SourceTreeImpl) GetBaseDir() string {
	return st.baseDir
}

// Watch starts watching for changes on a list of directories. The files under watch are the ones
// satisfying filenameMatcher() function. The watch lasts until either a watched file has been modified,
// or if an interruption is sent through the interrupt channel
func (st *SourceTreeImpl) Watch(
	dirList []string,
	filenameMatcher func(filename string) bool,
	interrupt <-chan bool,
) bool {
	// The file watcher
	st.watcher, _ = fsnotify.NewWatcher()
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			report.PostError(err)
		}
	}(st.watcher)

	// The filename matcher ensures that we watch only interesting files
	st.matcher = filenameMatcher

	// Used to notify if changes were detected on relevant files
	changesDetected := make(chan bool)

	// We recursively watch all subdirectories for all the provided directories
	for _, dir := range dirList {
		report.PostText("- watching ", dir)
		_ = filepath.Walk(dir, st.watchFile)
	}

	// Event handling goroutine
	go func() {
		for {
			select {
			case event := <-st.watcher.Events:
				report.PostText("-> ", event.Name)
				changesDetected <- true
				return
			case err := <-st.watcher.Errors:
				report.PostWarning(err)
				changesDetected <- false
				return
			case <-interrupt:
				changesDetected <- false
				return
			}
		}
	}()

	return <-changesDetected
}

// watchFile gets run as a walk func, searching for files to watch
func (st *SourceTreeImpl) watchFile(path string, fi os.FileInfo, err error) error {
	if err != nil {
		report.PostWarning(err)
		return err
	}

	// We don't watch directories themselves
	if fi.IsDir() {
		return nil
	}

	// If the filename matches our filter, we add it to the watching list
	if st.matcher(path) {
		err2 := st.watcher.Add(path)
		if err2 != nil {
			report.PostError(err2)
		}
		return err2
	}
	return nil
}
