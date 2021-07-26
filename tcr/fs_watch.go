package tcr

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mengdaming/tcr/trace"
	"os"
	"path/filepath"
)

var watcher *fsnotify.Watcher
var matcher func(filename string) bool

func WatchRecursive(
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

	// The filename matcher ensures that we watch only interesting files
	matcher = filenameMatcher

	// Used to notify if changes were detected on relevant files
	changesDetected := make(chan bool)

	// We recursively watch all subdirectories for all the provided directories
	for _, dir := range dirList {
		trace.Echo("- Watching ", dir)
		if err := filepath.Walk(dir, watchFile); err != nil {
			trace.Warning("filepath.Walk(", dir, "): ", err)
		}
	}

	// Event handling goroutine
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				trace.Echo("-> ", event.Name)
				changesDetected <- true
				return
			case err := <-watcher.Errors:
				trace.Warning("Watcher error: ", err)
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
func watchFile(path string, fi os.FileInfo, err error) error {
	if err != nil {
		trace.Warning("Something wrong with ", path)
		return err
	}

	// We don't watch directories themselves
	if fi.IsDir() {
		return nil
	}

	// If the filename matches our filter, we add it to the watching list
	if matcher(path) == true {
		err = watcher.Add(path)
		if err != nil {
			trace.Error("watcher.Add(", path, "): ", err)
		}
		return err
	}
	return nil
}
