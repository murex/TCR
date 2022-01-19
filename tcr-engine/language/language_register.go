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

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

var (
	builtIn    = make(map[string]LangInterface)
	registered = make(map[string]LangInterface)
)

// Register adds the provided language to the list of supported languages
func Register(lang LangInterface) error {
	if err := lang.checkName(); err != nil {
		return err
	}
	if err := lang.checkCompatibleToolchains(); err != nil {
		return err
	}
	if err := lang.checkDefaultToolchain(); err != nil {
		return err
	}
	registered[strings.ToLower(lang.GetName())] = lang
	return nil
}

func isSupported(name string) bool {
	_, found := registered[strings.ToLower(name)]
	return found
}

// Get returns the language instance with the provided name.
// The language name is case-insensitive.
// This method does not guarantee that the returned language instance can
// be used out of the box for file filtering operations as it does not
// enforce that baseDir is set. Prefer GetLanguage in this case.
func Get(name string) (LangInterface, error) {
	if name == "" {
		return nil, errors.New("language name not provided")
	}
	lang, found := registered[strings.ToLower(name)]
	if found {
		return lang, nil
	}
	return nil, errors.New(fmt.Sprint("language not supported: ", name))
}

// Names returns the list of available language names sorted alphabetically
func Names() []string {
	var names []string
	for _, lang := range registered {
		names = append(names, lang.GetName())
	}
	sort.Strings(names)
	return names
}

// Reset resets the language with the provided name to its default values
func Reset(name string) {
	_, found := registered[strings.ToLower(name)]
	if found && isBuiltIn(name) {
		_ = Register(getBuiltIn(name))
	}
}

func getBuiltIn(name string) LangInterface {
	var builtIn, _ = builtIn[strings.ToLower(name)]
	return builtIn
}

func isBuiltIn(name string) bool {
	_, found := builtIn[strings.ToLower(name)]
	return found
}

func addBuiltIn(lang LangInterface) error {
	if lang.GetName() == "" {
		return errors.New("language name cannot be an empty string")
	}
	builtIn[strings.ToLower(lang.GetName())] = lang
	return Register(lang)
}

// GetLanguage returns the language to be used in current session. If no value is provided
// for language (e.g. empty string), we try to detect the language based on the directory name.
// Both name and baseDir are case-insensitive
func GetLanguage(name string, baseDir string) (lang LangInterface, err error) {
	if name != "" {
		lang, err = getRegisteredLanguage(name)
	} else {
		lang, err = detectLanguageFromDirName(baseDir)
	}
	if lang != nil {
		lang.setBaseDir(baseDir)
	}
	return
}

func getRegisteredLanguage(name string) (LangInterface, error) {
	language, found := registered[strings.ToLower(name)]
	if found {
		return language, nil
	}
	return nil, errors.New(fmt.Sprint("language not supported: ", name))
}

// detectLanguageFromDirName is used to identify the language used in the provided directory. The current implementation
// simply looks at the name of the directory and checks if it matches with one of the supported languages
func detectLanguageFromDirName(baseDir string) (LangInterface, error) {
	return getRegisteredLanguage(filepath.Base(baseDir))
}
