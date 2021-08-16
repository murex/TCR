package language

import (
	"errors"
	"fmt"
	"path/filepath"
)

type Language interface {
	Name() string
	SrcDirs() []string
	TestDirs() []string
	IsSrcFile(filename string) bool
}

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

func DirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.SrcDirs(), lang.TestDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//report.PostInfo(dirList)
	return dirList
}
