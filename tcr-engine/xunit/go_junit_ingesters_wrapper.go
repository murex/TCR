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

package xunit

import (
	"github.com/joshdk/go-junit"
	"github.com/spf13/afero"
	"io"
	"os"
	"strings"
)

// The purpose of this file is to wrap go-junit Ingest functions so that they can work with afero in-memory file system.
// This file is a copy of
// (https://github.com/joshdk/go-junit/blob/6145f504ca0d053d5fd6c5379da09803fc3aef64/ingesters.go) where all
// non-afero-compliant calls are replaced by afero's ones. Go-junit Ingest functions are then called under the hood.
// Setting appFs variable to afero.NewMemMapFs() allows to use go-junit with afero's in-memory filesystem
// instead of the OS's filesystem.

var appFs afero.Fs

func init() {
	appFs = afero.NewOsFs()
}

// IngestDir will search the given directory for XML files and return a slice
// of all contained JUnit test suite definitions.
func IngestDir(directory string) ([]junit.Suite, error) {
	var filenames []string

	err := afero.Walk(appFs, directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Add all regular files that end with ".xml"
		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".xml") {
			filenames = append(filenames, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return IngestFiles(filenames)
}

// IngestFiles will parse the given XML files and return a slice of all
// contained JUnit test suite definitions.
func IngestFiles(filenames []string) ([]junit.Suite, error) {
	var all = make([]junit.Suite, 0)

	for _, filename := range filenames {
		suites, err := IngestFile(filename)
		if err != nil {
			return nil, err
		}
		all = append(all, suites...)
	}

	return all, nil
}

// IngestFile will parse the given XML file and return a slice of all contained
// JUnit test suite definitions.
func IngestFile(filename string) ([]junit.Suite, error) {
	file, err := appFs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file afero.File) {
		_ = file.Close()
	}(file)

	return IngestReader(file)
}

// IngestReader will parse the given XML reader and return a slice of all
// contained JUnit test suite definitions.
func IngestReader(reader io.Reader) ([]junit.Suite, error) {
	return junit.IngestReader(reader)
}

// Ingest will parse the given XML data and return a slice of all contained
// JUnit test suite definitions.
func Ingest(data []byte) ([]junit.Suite, error) {
	return junit.Ingest(data)
}
