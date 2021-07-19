package language

import (
	"path/filepath"
	"strings"
)

type Cpp struct {
}

func (language Cpp) Name() string {
	return "cpp"
}

func (language Cpp) SrcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

func (language Cpp) TestDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

func (language Cpp) IsSrcFile(filename string) bool {
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
