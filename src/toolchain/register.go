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

package toolchain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	builtIn    = make(map[string]TchnInterface)
	registered = make(map[string]TchnInterface)
)

// Register adds the provided toolchain to the list of supported toolchains
func Register(tchn TchnInterface) error {
	if err := tchn.checkName(); err != nil {
		return err
	}
	if err := tchn.checkBuildCommand(); err != nil {
		return err
	}
	if err := tchn.checkTestCommand(); err != nil {
		return err
	}
	registered[strings.ToLower(tchn.GetName())] = tchn
	return nil
}

func isSupported(name string) bool {
	_, found := registered[strings.ToLower(name)]
	return found
}

// Get returns the toolchain instance with the provided name
// The toolchain name is case insensitive.
func Get(name string) (TchnInterface, error) {
	if name == "" {
		return nil, errors.New("toolchain name not provided")
	}
	tchn, found := registered[strings.ToLower(name)]
	if found {
		return tchn, nil
	}
	return nil, errors.New(fmt.Sprint("toolchain not supported: ", name))
}

// Names returns the list of available toolchain names sorted alphabetically
func Names() []string {
	var names []string
	for _, tchn := range registered {
		names = append(names, tchn.GetName())
	}
	sort.Strings(names)
	return names
}

// Reset resets the toolchain with the provided name to its default values
func Reset(name string) {
	_, found := registered[strings.ToLower(name)]
	if found && isBuiltIn(name) {
		_ = Register(*getBuiltIn(name))
	}
}

// Unregister removes the toolchain with the provided name from the toolchain register
func Unregister(name string) {
	key := strings.ToLower(name)
	_, found := registered[key]
	if found {
		delete(registered, key)
	}
}

func getBuiltIn(name string) *TchnInterface {
	var builtIn, _ = builtIn[strings.ToLower(name)]
	return &builtIn
}

func isBuiltIn(name string) bool {
	_, found := builtIn[strings.ToLower(name)]
	return found
}

func addBuiltIn(tchn TchnInterface) error {
	if tchn.GetName() == "" {
		return errors.New("toolchain name cannot be an empty string")
	}
	builtIn[strings.ToLower(tchn.GetName())] = tchn
	return Register(tchn)
}
