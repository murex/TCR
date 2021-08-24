package language

import (
	"path/filepath"
	"strings"
)

// Java is the language implementation for java
type Java struct {
}

// Name returns the language name. This name is used to detect if a directory contains Java files
func (lang Java) Name() string {
	return "java"
}

// SrcDirs returns the list of subdirectories that may contain Java source files
func (lang Java) SrcDirs() []string {
	return []string{
		filepath.Join("src", "main"),
	}
}

// TestDirs returns the list of subdirectories that may contain Java test files
func (lang Java) TestDirs() []string {
	return []string{
		filepath.Join("src", "test"),
	}
}

// IsSrcFile returns true if the provided filename is recognized as a Java source file
func (lang Java) IsSrcFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".java":
		return true
	default:
		return false
	}
}
