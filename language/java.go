package language

import (
	"path/filepath"
	"strings"
)

type Java struct {
}

func (lang Java) Name() string {
	return "java"
}

func (lang Java) SrcDirs() []string {
	return []string{
		filepath.Join("src", "main"),
	}
}

func (lang Java) TestDirs() []string {
	return []string{
		filepath.Join("src", "test"),
	}
}

func (lang Java) IsSrcFile(filename string) bool {
	extension := filepath.Ext(filename)
	switch strings.ToLower(extension) {
	case ".java":
		return true
	default:
		return false
	}
}

