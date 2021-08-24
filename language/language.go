package language

import (
	"errors"
	"fmt"
	"path/filepath"
)

// Language is the interface that any supported language implementation must comply with
// in order to be used by TCR engine
type Language interface {
	Name() string
	SrcDirs() []string
	TestDirs() []string
	IsSrcFile(filename string) bool
}

// DetectLanguage is used to identify the language used in the provided directory. The current implementation
// simply looks at the name of the directory and checks if it matches with one of the supported languages
func DetectLanguage(baseDir string) (Language, error) {
	dir := filepath.Base(baseDir)
	switch dir {
	case "java":
		return Java{}, nil
	case "cpp":
		return Cpp{}, nil
	default:
		return nil, errors.New(fmt.Sprint("Unrecognized language: ", dir))
	}
}

// DirsToWatch returns the list of directories that TCR engine needs to watch for the provided language
func DirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.SrcDirs(), lang.TestDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//report.PostInfo(dirList)
	return dirList
}
