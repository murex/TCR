package tcr

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mengdaming/tcr/trace"
	"os"
	"path/filepath"
)

func watchFileSystem(
	dirList []string,
	filenameMatcher func(filename string) bool,
	interrupt <-chan bool,
) bool {
	trace.Info("Going to sleep until something interesting happens")
	return watchRecursive(dirList, filenameMatcher, interrupt)
}

var watcher *fsnotify.Watcher

func watchRecursive(
	dirList []string,
	filenameMatcher func(filename string) bool,
	interrupt <-chan bool,
) bool {

	// The file watcher
	watcher, _ = fsnotify.NewWatcher()
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			trace.Error("watcher.Close(): ", err)
		}
	}(watcher)

	// Used to indicate if changes were detected on relevant files
	changesDetected := make(chan bool)

	// We recursively watch all subdirectories for all the provided directories
	for _, dir := range dirList {
		trace.Transparent("- Watching ", dir)
		if err := filepath.Walk(dir, watchDir); err != nil {
			trace.Error("filepath.Walk(", dir, "): ", err)
		}
	}

	// Event handling goroutine
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					trace.Warning("<-watcher.Events: ", event)
					changesDetected <- false
					return
				}
				//trace.Info("Event:", event)
				if filenameMatcher(event.Name) {
					trace.Transparent("-> ", event.Name)
					changesDetected <- true
				} else {
					trace.Transparent("File change ignored: ", event.Name)
					changesDetected <- false
				}
				return
			case err, ok := <-watcher.Errors:
				trace.Warning("Watcher error: ", err)
				if !ok {
					changesDetected <- false
					return
				}
			case <-interrupt:
				//trace.Info("Ok, I stop watching")
				changesDetected <- false
				return
			}
		}
	}()

	return <-changesDetected
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func watchDir(path string, fi os.FileInfo, err error) error {

	if err != nil {
		trace.Warning("Something wrong with ", path)
		return err
	}
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		err = watcher.Add(path)
		if err != nil {
			trace.Error("watcher.Add(", path, "): ", err)
		}
		return err
	}
	return nil
}
