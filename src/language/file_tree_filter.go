/*
Copyright (c) 2023 Murex

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package language

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/murex/tcr/helpers"
	"github.com/spf13/afero"
)

// UnreachableDirectoryError is used to indicate that one or more directories cannot be accessed.
type UnreachableDirectoryError struct {
	dirs []string
}

// Error returns the error description
func (e *UnreachableDirectoryError) Error() string {
	var sb strings.Builder
	for _, dir := range e.DirList() {
		_, _ = sb.WriteString("cannot access directory: ")
		_, _ = sb.WriteString(dir)
		_, _ = sb.WriteRune('\n')
	}
	return sb.String()
}

// DirList returns the list of missing directories
func (e *UnreachableDirectoryError) DirList() []string {
	return e.dirs
}

// Add adds directories to the list of unreachable directories
func (e *UnreachableDirectoryError) Add(dir ...string) {
	e.dirs = append(e.dirs, dir...)
}

// IsEmpty indicates if the list of unreachable directories is empty
func (e *UnreachableDirectoryError) IsEmpty() bool {
	return len(e.dirs) == 0
}

type (
	// FileTreeFilter provides filtering mechanisms allowing to determine if a file or directory
	// is related to a language
	FileTreeFilter struct {
		Directories  []string
		FilePatterns []string
	}
)

func toLocalPath(input string) string {
	return filepath.Join(strings.Split(toSlashedPath(input), "/")...)
}

func toSlashedPath(input string) string {
	return path.Join(strings.Split(input, "\\")...)
}

func (ftf FileTreeFilter) isInFileTree(aPath string, baseDir string) bool {
	absPath, _ := filepath.Abs(aPath)
	// If no directory is configured, any path that is under baseDir path is ok
	if len(ftf.Directories) == 0 {
		if helpers.IsSubPathOf(absPath, baseDir) {
			return true
		}
	}

	for _, dir := range ftf.Directories {
		filterAbsPath, _ := filepath.Abs(filepath.Join(baseDir, dir))
		if helpers.IsSubPathOf(absPath, filterAbsPath) {
			return true
		}
	}
	return false
}

func (ftf FileTreeFilter) matches(p string, baseDir string) bool {
	if p == "" {
		return false
	}
	if ftf.isInFileTree(p, baseDir) {
		// If no pattern is set, any file matches as long as it's in the file tree
		if len(ftf.FilePatterns) == 0 {
			return true
		}
		for _, filter := range ftf.FilePatterns {
			re := regexp.MustCompile(filter)
			if re.MatchString(p) {
				return true
			}
		}
	}
	return false
}

func (ftf FileTreeFilter) findAllMatchingFiles(baseDir string) (files []string, err error) {
	dirErr := UnreachableDirectoryError{}

	for _, dir := range ftf.Directories {
		matchingFiles, err := ftf.findMatchingFilesInDir(baseDir, dir)
		if err != nil {
			dirErr.Add(filepath.Join(baseDir, dir))
			continue
		}
		files = append(files, matchingFiles...)
	}

	if !dirErr.IsEmpty() {
		return files, &dirErr
	}
	return files, nil
}

func (ftf FileTreeFilter) findMatchingFilesInDir(baseDir string, dir string) (files []string, err error) {
	return files, afero.Walk(appFS, filepath.Join(baseDir, dir),
		func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}
			if ftf.matches(path, baseDir) {
				files = append(files, path)
				return err
			}
			return nil
		})
}
