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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

const (
	baseDir   = "base-dir"
	baseName  = "a-name"
	extension = ".yml"
)

func Test_create_sub_dir(t *testing.T) {
	const existingFile = "already-used"
	tests := []struct {
		desc          string
		dirPath       string
		expectedTrace []string
	}{
		{
			"empty dir path",
			"",
			[]string{},
		},
		{
			"existing path",
			baseDir,
			[]string{},
		},
		{
			"non-existing path",
			"xxx",
			[]string{"Creating non-existing path: xxx"},
		},
		{
			"non-existing sub-path",
			path.Join(baseDir, "xxx"),
			[]string{"Creating non-existing sub-path: base-dir/xxx"},
		},
		{
			"path already used",
			path.Join(baseDir, existingFile),
			[]string{"Error creating path already used: a file with this name already exists"},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_ = fs.MkdirAll(baseDir, 0755)
			_, _ = fs.Create(filepath.Join(baseDir, existingFile))
			AssertSimpleTrace(t, test.expectedTrace, func() {
				CreateSubDir(fs, test.dirPath, test.desc)
			})
		})
	}
}

func Test_build_yaml_filepath(t *testing.T) {
	dirPath := "some-dir/some-sub-dir"
	assert.Equal(t, filepath.Join(dirPath, baseName+extension), BuildYAMLFilePath(dirPath, baseName))
}

func Test_can_retrieve_yaml_filename_from_name(t *testing.T) {
	assert.Equal(t, baseName+extension, buildYAMLFilename(baseName))
}

func Test_yaml_filename_is_always_lowercase(t *testing.T) {
	assert.Equal(t, baseName+extension, buildYAMLFilename(strings.ToUpper(baseName)))
}

func Test_can_retrieve_name_from_yaml_filename(t *testing.T) {
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToLower(baseName)+strings.ToLower(extension)))
}

func Test_name_from_yaml_filename_is_always_lowercase(t *testing.T) {
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToUpper(baseName)+strings.ToLower(extension)))
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToLower(baseName)+strings.ToUpper(extension)))
	assert.Equal(t, baseName, ExtractNameFromYAMLFilename(strings.ToUpper(baseName)+strings.ToUpper(extension)))
}

func Test_list_yaml_files_in(t *testing.T) {
	yaml1 := "file1" + extension
	yaml2 := "file2" + extension
	other := "other.txt"

	tests := []struct {
		desc     string
		dir      string
		files    []string
		expected []string
	}{
		{
			"empty dir",
			"base-dir",
			[]string{},
			nil,
		},
		{
			"0 yaml 1 other",
			"base-dir",
			[]string{other},
			nil,
		},
		{
			"1 yaml 0 other",
			"base-dir",
			[]string{yaml1},
			[]string{yaml1},
		},
		{
			"2 yaml 1 other",
			"base-dir",
			[]string{yaml1, yaml2, other},
			[]string{yaml1, yaml2},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			_ = fs.MkdirAll(test.dir, 0755)
			for _, file := range test.files {
				_, _ = fs.Create(filepath.Join(test.dir, file))
			}
			assert.Equal(t, test.expected, ListYAMLFilesIn(fs, test.dir))
		})
	}
}

type testStructYAML struct {
	Field1 string `yaml:"field1"`
	Field2 string `yaml:"field2"`
}

const yamlData = "field1: value1\nfield2: value2\n"

func Test_load_from_yaml_file(t *testing.T) {
	const yamlFile = "test.yml"
	dir, errTempDir := os.MkdirTemp("", "tcr-utils")
	if errTempDir != nil {
		t.Fatal(errTempDir)
	}
	defer t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	yamlFilePath := filepath.Join(dir, yamlFile)
	errCreate := os.WriteFile(yamlFilePath, []byte(yamlData), 0600)
	if errCreate != nil {
		t.Fatal(errCreate)
	}
	var yaml testStructYAML
	errLoad := LoadFromYAMLFile(os.DirFS(dir), yamlFile, &yaml)
	if assert.NoError(t, errLoad) {
		assert.Equal(t, "value1", yaml.Field1)
		assert.Equal(t, "value2", yaml.Field2)
	}
}

func Test_save_to_yaml_file(t *testing.T) {
	yaml := testStructYAML{Field1: "value1", Field2: "value2"}
	const yamlFile = "test.yml"

	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll(baseDir, 0755)
	yamlFilePath := filepath.Join(baseDir, yamlFile)
	AssertSimpleTrace(t, nil, func() {
		// When there is an error, trace is not empty
		SaveToYAMLFile(fs, &yaml, yamlFilePath)
	})
	// Read file contents and compare it to what we expect
	fileContents, err := afero.ReadFile(fs, yamlFilePath)
	if assert.NoError(t, err) {
		assert.Equal(t, yamlData, string(fileContents))
	}
}
