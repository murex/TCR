package language

import (
	"path/filepath"
	"strings"
)

type Cpp struct {
}

func (lang Cpp) Name() string {
	return "cpp"
}

func (lang Cpp) SrcDirs() []string {
	return []string{
		filepath.Join("src"),
		filepath.Join("include"),
	}
}

func (lang Cpp) TestDirs() []string {
	return []string{
		filepath.Join("test"),
	}
}

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
