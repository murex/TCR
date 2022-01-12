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

type (
	// FakeLanguage is a Language fake that makes sure there is no filesystem
	// access. It overwrites AllSrcFiles() so that it always returns a non-empty
	// list of filenames
	FakeLanguage struct {
		Language
	}
)

// NewFakeLanguage creates a new Language instance compatible with the provided toolchain
func NewFakeLanguage(toolchainName string) *FakeLanguage {
	return &FakeLanguage{
		Language: Language{
			name: "fake-language",
			toolchains: Toolchains{
				Default:    toolchainName,
				Compatible: []string{toolchainName},
			},
		},
	}
}

// AllSrcFiles returns the list of source files for this language.
// Always returns a list of 2 fake filenames.
func (lang *FakeLanguage) AllSrcFiles() (result []string, err error) {
	result = []string{"fake-file1", "fake-file2"}
	return
}
