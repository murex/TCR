package filesystem

import (
	"github.com/mengdaming/tcr/trace"

	"github.com/fsnotify/fsnotify"

	"os"
	"path/filepath"
)

type SourceTreeImpl struct {
	baseDir string
	watcher *fsnotify.Watcher
	matcher func(filename string) bool
}

func NewSourceTreeImpl(dir string) SourceTree {
	var st = SourceTreeImpl{}
	var err error
	st.baseDir, err = st.changeDir(dir)
	if err != nil {
		return nil
	}
	return &st
}

func (st *SourceTreeImpl) changeDir(dir string) (string, error) {
	_, err := os.Stat(dir)
	switch {
	case os.IsNotExist(err):
		trace.Error("Directory ", dir, " does not exist")
		return "", err
	case os.IsPermission(err):
		trace.Error("Can't access directory ", dir)
		return "", err
	}

	err = os.Chdir(dir)
	if err != nil {
		trace.Error("Failed to change directory to ", dir)
		return "", err
	}

	wd, _ := os.Getwd()
	return wd, nil
}

func (st *SourceTreeImpl) GetBaseDir() string {
	return st.baseDir
}

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
			trace.Error("watcher.Close(): ", err)
		}
	}(st.watcher)

	// The filename matcher ensures that we watch only interesting files
	st.matcher = filenameMatcher

	// Used to notify if changes were detected on relevant files
	changesDetected := make(chan bool)

	// We recursively watch all subdirectories for all the provided directories
	for _, dir := range dirList {
		trace.Echo("- Watching ", dir)
		if err := filepath.Walk(dir, st.watchFile); err != nil {
			trace.Warning("filepath.Walk(", dir, "): ", err)
		}
	}

	// Event handling goroutine
	go func() {
		for {
			select {
			case event := <-st.watcher.Events:
				trace.Echo("-> ", event.Name)
				changesDetected <- true
				return
			case err := <-st.watcher.Errors:
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
func (st *SourceTreeImpl) watchFile(path string, fi os.FileInfo, err error) error {
	if err != nil {
		trace.Warning("Something wrong with ", path)
		return err
	}

	// We don't watch directories themselves
	if fi.IsDir() {
		return nil
	}

	// If the filename matches our filter, we add it to the watching list
	if st.matcher(path) == true {
		err = st.watcher.Add(path)
		if err != nil {
			trace.Error("watcher.Add(", path, "): ", err)
		}
		return err
	}
	return nil
}
