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

package utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	baseName  = "a-name"
	extension = ".yml"
)

func Test_can_retrieve_name_from_yaml_filename(t *testing.T) {
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToLower(baseName)+strings.ToLower(extension)))
}

func Test_name_from_yaml_filename_is_always_lowercase(t *testing.T) {
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToUpper(baseName)+strings.ToLower(extension)))
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToLower(baseName)+strings.ToUpper(extension)))
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToUpper(baseName)+strings.ToUpper(extension)))
}

func Test_can_retrieve_yaml_filename_from_name(t *testing.T) {
	assert.Equal(t, baseName+extension, buildYAMLFilename(baseName))
}

func Test_yaml_filename_is_always_lowercase(t *testing.T) {
	assert.Equal(t, baseName+extension, buildYAMLFilename(strings.ToUpper(baseName)))
}
