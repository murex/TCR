package tcr

import "path"

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
	return []string{path.Join(BaseDir(), "src", "main")}
}

func (Language JavaLanguage) testDirs() []string {
	return []string{path.Join(BaseDir(), "src", "test")}
}