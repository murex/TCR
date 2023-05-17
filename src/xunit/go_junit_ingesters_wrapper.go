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
	"github.com/mengdaming/go-junit"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// The purpose of this file is to wrap go-junit ingest functions so that they can work with afero in-memory file system.
// This file is a copy of
// (https://github.com/joshdk/go-junit/blob/6145f504ca0d053d5fd6c5379da09803fc3aef64/ingesters.go) where all
// non-afero-compliant calls are replaced by afero's ones. Go-junit ingest functions are then called under the hood.
// Setting appFs variable to afero.NewMemMapFs() allows to use go-junit with afero's in-memory filesystem
// instead of the OS's filesystem.

var appFs afero.Fs

func init() {
	appFs = afero.NewOsFs()
}

// ingestDir will search the given directory for XML files and return a slice
// of all contained JUnit test suite definitions.
func ingestDir(directory string) ([]junit.Suite, error) {
	var filenames []string

	d, errSymLink := evalSymLink(directory)
	if errSymLink != nil {
		return nil, errSymLink
	}

	errWalk := afero.Walk(appFs, d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Add all regular files that end with ".xml"
		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".xml") {
			filenames = append(filenames, path)
		}
		return nil
	})
	if errWalk != nil {
		return nil, errWalk
	}

	return ingestFiles(filenames)
}

// evalSymLink tries to convert a symbolic link path to the path it points to.
// Warning: afero.MemMapFs does not support symbolic links. For this reason,
// symbolic-link related tests need to be run with real OS filesystem.
// Cf. https://github.com/spf13/afero/issues/258
// If the provided path is a regular path, the function returns it unchanged
func evalSymLink(directory string) (string, error) {
	// afero.MemMapFs does not support symbolic links
	_, symLinkSupported, errLstat := appFs.(afero.Lstater).LstatIfPossible(directory)
	if errLstat != nil {
		return "", errLstat
	}
	if !symLinkSupported {
		return directory, nil
	}
	// afero.Walk does not follow symbolic links.
	// So we make sure to resolve the parent directory before calling afero.Walk()
	d, errSymLink := filepath.EvalSymlinks(directory)
	if errSymLink != nil {
		return "", errSymLink
	}
	return d, nil
}

// ingestFiles will parse the given XML files and return a slice of all
// contained JUnit test suite definitions.
func ingestFiles(filenames []string) ([]junit.Suite, error) {
	var all = make([]junit.Suite, 0)

	for _, filename := range filenames {
		suites, err := ingestFile(filename)
		if err != nil {
			return nil, err
		}
		all = append(all, suites...)
	}

	return all, nil
}

// ingestFile will parse the given XML file and return a slice of all contained
// JUnit test suite definitions.
func ingestFile(filename string) ([]junit.Suite, error) {
	file, err := appFs.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file afero.File) {
		_ = file.Close()
	}(file)

	return ingestReader(file)
}

// ingestReader will parse the given XML reader and return a slice of all
// contained JUnit test suite definitions.
func ingestReader(reader io.Reader) ([]junit.Suite, error) {
	return junit.IngestReader(reader)
}

// ingest will parse the given XML data and return a slice of all contained
// JUnit test suite definitions.
func ingest(data []byte) ([]junit.Suite, error) {
	return junit.Ingest(data)
}
