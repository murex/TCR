package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"path/filepath"
	"strings"
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
		return JavaLanguage{}
	case "cpp":
		return CppLanguage{}
	default:
		trace.Error("Unrecognized language: ", dir)
	}
	return nil
}

// ========================================================================

type JavaLanguage struct {
}

func (language JavaLanguage) Name() string {
	return "java"
}

func (language JavaLanguage) SrcDirs() []string {
	return []string{
		filepath.Join("src", "main"),
	}
}

func (language JavaLanguage) TestDirs() []string {
	return []string{
		filepath.Join("src", "test"),
	}
}

func (language JavaLanguage) IsSrcFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".java":
		return true
	default:
		return false
	}
}

// ========================================================================

type CppLanguage struct {
}

func (language CppLanguage) Name() string {
	return "cpp"
}

func (language CppLanguage) SrcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

func (language CppLanguage) TestDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

func (language CppLanguage) IsSrcFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".cpp", ".hpp":
		return true
	case ".c", ".h":
		return true
	case ".cc", ".hh":
		return true
	default:
		return false
	}
}

// ========================================================================

func DirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.SrcDirs(), lang.TestDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//trace.Info(dirList)
	return dirList
}
