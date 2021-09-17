package language

import (
	"path/filepath"
	"strings"
)

// Cpp is the language implementation for C++
type Cpp struct {
}

// Name returns the language name. This name is used to detect if a directory contains Cpp files
func (lang Cpp) Name() string {
	return "cpp"
}

// SrcDirs returns the list of subdirectories that may contain Cpp source files
func (lang Cpp) SrcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

// TestDirs returns the list of subdirectories that may contain Cpp test files
func (lang Cpp) TestDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

// IsSrcFile returns true if the provided filename is recognized as a Cpp source file
func (lang Cpp) IsSrcFile(filename string) bool {
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
