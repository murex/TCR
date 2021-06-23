package tcr

import (
	"path/filepath"
)

type Language interface {
	name() string
	toolchain() string
	workDir() string
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

func (Language JavaLanguage) workDir() string {
	return BaseDir()
}

func (Language JavaLanguage) srcDirs() []string {
	return []string{
		filepath.Join(BaseDir(), "src", "main"),
	}
}

func (Language JavaLanguage) testDirs() []string {
	return []string{
		filepath.Join(BaseDir(), "src", "test"),
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

func (Language CppLanguage) workDir() string {
	return filepath.Join(BaseDir(), "build")
}

func (Language CppLanguage) srcDirs() []string {
	return []string{
		filepath.Join(BaseDir(), "src"),
		filepath.Join(BaseDir(), "include"),
	}
}

func (Language CppLanguage) testDirs() []string {
	return []string{
		filepath.Join(BaseDir(), "test"),
	}
}
