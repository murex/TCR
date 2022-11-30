//go:build test_helper

/*
Copyright (c) 2022 Murex

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

// AFileTreeFilter is a test data builder for type FileTreeFilter
func AFileTreeFilter(builders ...func(filter *FileTreeFilter)) *FileTreeFilter {
	filter := &FileTreeFilter{
		Directories:  []string{},
		FilePatterns: []string{},
	}

	for _, build := range builders {
		build(filter)
	}
	return filter
}

// WithDirectory adds the provided dirName to the list of directories covered by the created FileTreeFilter
func WithDirectory(dirName string) func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.Directories = append(filter.Directories, dirName)
	}
}

// WithDirectories sets the list of directories covered by the created FileTreeFilter
func WithDirectories(dirs ...string) func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.Directories = dirs
	}
}

// WithNoDirectory enforces that the created FileTreeFilter has no directory configured
func WithNoDirectory() func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.Directories = nil
	}
}

// WithPattern adds the provided pattern to the list of filename patterns recognized by the created FileTreeFilter
func WithPattern(pattern string) func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.FilePatterns = append(filter.FilePatterns, pattern)
	}
}

// WithPatterns sets the list of filename patterns recognized by the created FileTreeFilter
func WithPatterns(patterns ...string) func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.FilePatterns = patterns
	}
}

// WithOpenPattern adds a pattern matching any filename to the created FileTreeFilter
func WithOpenPattern() func(filter *FileTreeFilter) {
	return WithPattern("^.*$")
}

// WithClosedPattern adds a pattern not matching any filename to the created FileTreeFilter
func WithClosedPattern() func(filter *FileTreeFilter) {
	return WithPattern("^$")
}

// WithNoPattern enforces that the created FileTreeFilter has no filename pattern
func WithNoPattern() func(filter *FileTreeFilter) {
	return func(filter *FileTreeFilter) {
		filter.FilePatterns = nil
	}
}
