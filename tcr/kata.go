package tcr

import (
	"github.com/mengdaming/tcr/trace"
	"os"
	"path"
)

type Kata struct {
}

func BaseDir() string {
	baseDir, err := os.Executable()
	if err != nil {
		trace.Error(err.Error())
	}
	return baseDir
}

func ScriptDir() string {
	return path.Join(BaseDir(), "tcr")
}