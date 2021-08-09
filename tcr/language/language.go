package language

import (
	"github.com/mengdaming/tcr/tcr/report"
	"path/filepath"
)

type Language interface {
	Name() string
	SrcDirs() []string
	TestDirs() []string
	IsSrcFile(filename string) bool
}

func DetectLanguage(baseDir string) Language {
	dir := filepath.Base(baseDir)
	switch dir {
	case "java":
		return Java{}
	case "cpp":
		return Cpp{}
	default:
		report.PostError("Unrecognized language: ", dir)
	}
	return nil
}

func DirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.SrcDirs(), lang.TestDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//report.PostInfo(dirList)
	return dirList
}
