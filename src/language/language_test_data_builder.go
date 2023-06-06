//go:build test_helper

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

// ALanguage is a test data builder for type Language
func ALanguage(languageBuilders ...func(lang *Language)) *Language {
	lang := New("default-language",
		Toolchains{Default: "default-toolchain", Compatible: []string{"default-toolchain"}},
		*AFileTreeFilter(),
		*AFileTreeFilter(),
	)

	for _, build := range languageBuilders {
		build(lang)
	}
	return lang
}

// WithName allows to create a language with the provided name
func WithName(name string) func(lang *Language) {
	return func(lang *Language) { lang.name = name }
}

// WithNoCompatibleToolchain allows to create a language with no compatible toolchain defined
func WithNoCompatibleToolchain() func(lang *Language) {
	return func(lang *Language) { lang.toolchains.Compatible = nil }
}

// WithCompatibleToolchain adds the provided toolchain to the list of compatible toolchains defined for this language
func WithCompatibleToolchain(tchn string) func(lang *Language) {
	return func(lang *Language) {
		lang.toolchains.Compatible = append(lang.toolchains.Compatible, tchn)
	}
}

// WithNoDefaultToolchain allows to create a language with no default toolchain defined
func WithNoDefaultToolchain() func(lang *Language) {
	return func(lang *Language) { lang.toolchains.Default = "" }
}

// WithDefaultToolchain sets the provided toolchain as the default toolchain for this language
func WithDefaultToolchain(tchn string) func(lang *Language) {
	return func(lang *Language) { lang.toolchains.Default = tchn }
}

// WithSrcFiles sets the provided filter for source files for this language
func WithSrcFiles(filter *FileTreeFilter) func(lang *Language) {
	return func(lang *Language) { lang.srcFileFilter = *filter }
}

// WithTestFiles sets the provided filter for test files for this language
func WithTestFiles(filter *FileTreeFilter) func(lang *Language) {
	return func(lang *Language) { lang.testFileFilter = *filter }
}

// WithBaseDir sets the provided directory as base directory for this language
func WithBaseDir(dir string) func(lang *Language) {
	return func(lang *Language) { lang.baseDir = dir }
}
