package tcr

import (
	"path/filepath"
)

type Language interface {
	name() string
	toolchain() string
	srcDirs() []string
	testDirs() []string
}

// ========================================================================

type JavaLanguage struct {
}

func (Language JavaLanguage) name() string {
	return "java"
}

func (Language JavaLanguage) toolchain() string {
	return "gradle"
}

func (Language JavaLanguage) srcDirs() []string {
	return []string{
		filepath.Join("src", "main"),
	}
}

func (Language JavaLanguage) testDirs() []string {
	return []string{
		filepath.Join("src", "test"),
	}
}

// ========================================================================

type CppLanguage struct {
}

func (Language CppLanguage) name() string {
	return "cpp"
}

func (Language CppLanguage) toolchain() string {
	return "cmake"
}

func (Language CppLanguage) srcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

func (Language CppLanguage) testDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

// ========================================================================

func dirsToWatch(baseDir string, lang Language) []string {
	dirList := append(lang.srcDirs(), lang.testDirs()...)
	for i := 0; i < len(dirList); i++ {
		dirList[i] = filepath.Join(baseDir, dirList[i])
	}
	//trace.Info(dirList)
	return dirList
}
