/*
Copyright (c) 2021 Murex

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
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

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

func (treeFilter FileTreeFilter) isInFileTree(path string, baseDir string) bool {
	absPath, _ := filepath.Abs(path)
	// If no directory is configured, any path that is under baseDir path is ok
	if treeFilter.Directories == nil || len(treeFilter.Directories) == 0 {
		if isSubPathOf(absPath, baseDir) {
			return true
		}
	}

	for _, dir := range treeFilter.Directories {
		filterAbsPath, _ := filepath.Abs(filepath.Join(baseDir, dir))
		if isSubPathOf(absPath, filterAbsPath) {
			return true
		}
	}
	return false
}

func isSubPathOf(path string, refPath string) bool {
	// If refPath is empty, we consider it as being the root, thus path is a sub-path of refPath
	if refPath == "" {
		return true
	}
	if refPath == path || strings.HasPrefix(path, refPath+string(os.PathSeparator)) {
		return true
	}
	return false
}

func (treeFilter FileTreeFilter) matches(filepath string, baseDir string) bool {
	if filepath == "" {
		return false
	}
	if treeFilter.isInFileTree(filepath, baseDir) {
		// If no pattern is set, any file matches as long as it's in the file tree
		if treeFilter.FilePatterns == nil || len(treeFilter.FilePatterns) == 0 {
			return true
		}
		for _, filter := range treeFilter.FilePatterns {
			re := regexp.MustCompile(filter)
			if re.MatchString(filepath) {
				return true
			}
		}
	}
	return false
}

// buildRegex is a utility function adding flags and markers to a regex pattern
// so that it ignores character case and has the beginning and end of text markers
func buildRegex(corePattern string) string {
	const (
		ignoreCaseFlag  = "(?i)"
		beginningMarker = "^"
		endMarker       = "$"
	)
	return ignoreCaseFlag + beginningMarker + corePattern + endMarker
}

func (treeFilter FileTreeFilter) findAllMatchingFiles(baseDir string) (files []string, err error) {
	for _, dir := range treeFilter.Directories {
		err := filepath.Walk(filepath.Join(baseDir, dir), func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return errors.New("something wrong with " + path)
			}
			if fi.IsDir() {
				return nil
			}
			// If the filename matches the file pattern, we add it to the list of files
			if treeFilter.matches(path, baseDir) {
				files = append(files, path)
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
