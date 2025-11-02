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

package checker

import (
	"errors"
	"testing"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/filesystem"
	"github.com/murex/tcr/language"
	"github.com/murex/tcr/params"
	"github.com/stretchr/testify/assert"
)

func Test_check_language(t *testing.T) {
	assertCheckGroupRunner(t,
		checkLanguage,
		&checkLanguageRunners,
		*params.AParamSet(),
		"language")
}

func Test_check_language_parameter(t *testing.T) {
	tests := []struct {
		desc     string
		lang     string
		langErr  error
		expected []model.CheckPoint
	}{
		{"not set", "", nil, nil},
		{
			"set and valid", "valid-language", nil,
			[]model.CheckPoint{
				model.OkCheckPoint("language parameter is set to valid-language"),
				model.OkCheckPoint("valid-language language is valid"),
			},
		},
		{
			"set but invalid", "wrong-language", errors.New("language parameter error"),
			[]model.CheckPoint{
				model.OkCheckPoint("language parameter is set to wrong-language"),
				model.ErrorCheckPoint("language parameter error"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.lang))
			checkEnv.langErr = test.langErr
			assert.Equal(t, test.expected, checkLanguageParameter(p))
		})
	}
}

func Test_check_language_detection(t *testing.T) {
	tests := []struct {
		desc         string
		langParam    string
		langDetected language.LangInterface
		langErr      error
		sourceTree   filesystem.SourceTree
		expected     []model.CheckPoint
	}{
		{"language param is set",
			"some-language", nil, nil, nil, nil},
		{
			"invalid source tree",
			"", nil, nil,
			nil,
			[]model.CheckPoint{
				model.OkCheckPoint("language parameter is not set explicitly"),
				model.ErrorCheckPoint("cannot retrieve language from base directory name"),
			},
		},
		{
			"invalid language",
			"", nil, errors.New("wrong language"),
			filesystem.NewFakeSourceTree("/some-path/wrong-language"),
			[]model.CheckPoint{
				model.OkCheckPoint("language parameter is not set explicitly"),
				model.OkCheckPoint("using base directory name for language detection"),
				model.OkCheckPoint("base directory is /some-path/wrong-language"),
				model.ErrorCheckPoint("wrong language"),
			},
		},
		{
			"all green",
			"", language.ALanguage(language.WithName("java")), nil,
			filesystem.NewFakeSourceTree("/some-path/java"),
			[]model.CheckPoint{
				model.OkCheckPoint("language parameter is not set explicitly"),
				model.OkCheckPoint("using base directory name for language detection"),
				model.OkCheckPoint("base directory is /some-path/java"),
				model.OkCheckPoint("language retrieved from base directory name: java"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langParam))
			checkEnv.sourceTree = test.sourceTree
			checkEnv.lang = test.langDetected
			checkEnv.langErr = test.langErr
			assert.Equal(t, test.expected, checkLanguageDetection(p))
		})
	}
}

func Test_check_language_src_directories(t *testing.T) {
	tests := []struct {
		desc     string
		langName string
		lang     language.LangInterface
		expected []model.CheckPoint
	}{
		{"no language", "", nil, nil},
		{
			"0 src dir", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithSrcFiles(
					language.AFileTreeFilter(language.WithNoDirectory()))),
			[]model.CheckPoint{
				model.WarningCheckPoint("no source directory defined for xxx language"),
			},
		},
		{
			"1 src dir", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithSrcFiles(
					language.AFileTreeFilter(language.WithDirectory("dir1")))),
			[]model.CheckPoint{
				model.OkCheckPoint("source directories:"),
				model.OkCheckPoint("- dir1"),
			},
		},
		{
			"2 src dirs", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithSrcFiles(
					language.AFileTreeFilter(language.WithDirectories("dir1", "dir2")))),
			[]model.CheckPoint{
				model.OkCheckPoint("source directories:"),
				model.OkCheckPoint("- dir1"),
				model.OkCheckPoint("- dir2"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langName))
			checkEnv.lang = test.lang
			assert.Equal(t, test.expected, checkLanguageSrcDirectories(p))
		})
	}
}

func Test_check_language_src_patterns(t *testing.T) {
	tests := []struct {
		desc     string
		langName string
		lang     language.LangInterface
		expected []model.CheckPoint
	}{
		{"no language", "", nil, nil},
		{
			"0 src pattern", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithSrcFiles(
					language.AFileTreeFilter(language.WithNoPattern()))),
			[]model.CheckPoint{
				model.WarningCheckPoint("no source filename pattern defined for xxx language"),
			},
		},
		{
			"1 src pattern", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithSrcFiles(
					language.AFileTreeFilter(language.WithPattern("*.xxx")))),
			[]model.CheckPoint{
				model.OkCheckPoint("source filename matching patterns:"),
				model.OkCheckPoint("- *.xxx"),
			},
		},
		{
			"2 src patterns", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithSrcFiles(
					language.AFileTreeFilter(language.WithPatterns("*.xxx", "*.yyy")))),
			[]model.CheckPoint{
				model.OkCheckPoint("source filename matching patterns:"),
				model.OkCheckPoint("- *.xxx"),
				model.OkCheckPoint("- *.yyy"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langName))
			checkEnv.lang = test.lang
			assert.Equal(t, test.expected, checkLanguageSrcPatterns(p))
		})
	}
}

func Test_check_language_src_files(t *testing.T) {
	tests := []struct {
		desc     string
		langName string
		lang     language.LangInterface
		expected []model.CheckPoint
	}{
		{"no language", "", nil, nil},
		{
			"0 match", "xxx",
			language.NewFakeLanguage("").WithAllSrcFiles(
				func() ([]string, error) {
					return []string{}, nil
				}),
			[]model.CheckPoint{
				model.WarningCheckPoint("no matching source file found"),
			},
		},
		{
			"1 match", "xxx",
			language.NewFakeLanguage("").WithAllSrcFiles(
				func() ([]string, error) {
					return []string{"src-file1"}, nil
				}),
			[]model.CheckPoint{
				model.OkCheckPoint("matching source files found:"),
				model.OkCheckPoint("- src-file1"),
			},
		},
		{
			"2 matches", "xxx",
			language.NewFakeLanguage("").WithAllSrcFiles(
				func() ([]string, error) {
					return []string{"src-file1", "src-file2"}, nil
				}),
			[]model.CheckPoint{
				model.OkCheckPoint("matching source files found:"),
				model.OkCheckPoint("- src-file1"),
				model.OkCheckPoint("- src-file2"),
			},
		},
		{
			"matching error", "xxx",
			language.NewFakeLanguage("").WithAllSrcFiles(
				func() ([]string, error) {
					return nil, errors.New("some error")
				}),
			[]model.CheckPoint{
				model.ErrorCheckPoint("some error"),
			},
		},
		{
			"unreachable directory error and no matching file", "xxx",
			language.NewFakeLanguage("").WithAllSrcFiles(
				func() ([]string, error) {
					err := language.UnreachableDirectoryError{}
					err.Add("dir1")
					return nil, &err
				}),
			[]model.CheckPoint{
				model.WarningCheckPoint("cannot access source directory dir1"),
				model.WarningCheckPoint("no matching source file found"),
			},
		},
		{
			"unreachable directory error and one matching file", "xxx",
			language.NewFakeLanguage("").WithAllSrcFiles(
				func() ([]string, error) {
					err := language.UnreachableDirectoryError{}
					err.Add("dir1")
					return []string{"src-file1"}, &err
				}),
			[]model.CheckPoint{
				model.WarningCheckPoint("cannot access source directory dir1"),
				model.OkCheckPoint("matching source files found:"),
				model.OkCheckPoint("- src-file1"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langName))
			checkEnv.lang = test.lang
			assert.Equal(t, test.expected, checkLanguageSrcFiles(p))
		})
	}
}

func Test_check_language_test_directories(t *testing.T) {
	tests := []struct {
		desc     string
		langName string
		lang     language.LangInterface
		expected []model.CheckPoint
	}{
		{"no language", "", nil, nil},
		{
			"0 test dir", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithTestFiles(
					language.AFileTreeFilter(language.WithNoDirectory()))),
			[]model.CheckPoint{
				model.WarningCheckPoint("no test directory defined for xxx language"),
			},
		},
		{
			"1 test dir", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithTestFiles(
					language.AFileTreeFilter(language.WithDirectory("dir1")))),
			[]model.CheckPoint{
				model.OkCheckPoint("test directories:"),
				model.OkCheckPoint("- dir1"),
			},
		},
		{
			"2 test dirs", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithTestFiles(
					language.AFileTreeFilter(language.WithDirectories("dir1", "dir2")))),
			[]model.CheckPoint{
				model.OkCheckPoint("test directories:"),
				model.OkCheckPoint("- dir1"),
				model.OkCheckPoint("- dir2"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langName))
			checkEnv.lang = test.lang
			assert.Equal(t, test.expected, checkLanguageTestDirectories(p))
		})
	}
}

func Test_check_language_test_patterns(t *testing.T) {
	tests := []struct {
		desc     string
		langName string
		lang     language.LangInterface
		expected []model.CheckPoint
	}{
		{"no language", "", nil, nil},
		{
			"0 test pattern", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithTestFiles(
					language.AFileTreeFilter(language.WithNoPattern()))),
			[]model.CheckPoint{
				model.WarningCheckPoint("no test filename pattern defined for xxx language"),
			},
		},
		{
			"1 test pattern", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithTestFiles(
					language.AFileTreeFilter(language.WithPattern("*.xxx")))),
			[]model.CheckPoint{
				model.OkCheckPoint("test filename matching patterns:"),
				model.OkCheckPoint("- *.xxx"),
			},
		},
		{
			"2 test patterns", "xxx",
			language.ALanguage(language.WithName("xxx"),
				language.WithTestFiles(
					language.AFileTreeFilter(language.WithPatterns("*.xxx", "*.yyy")))),
			[]model.CheckPoint{
				model.OkCheckPoint("test filename matching patterns:"),
				model.OkCheckPoint("- *.xxx"),
				model.OkCheckPoint("- *.yyy"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langName))
			checkEnv.lang = test.lang
			assert.Equal(t, test.expected, checkLanguageTestPatterns(p))
		})
	}
}

func Test_check_language_test_files(t *testing.T) {
	tests := []struct {
		desc     string
		langName string
		lang     language.LangInterface
		expected []model.CheckPoint
	}{
		{"no language", "", nil, nil},
		{
			"0 match", "xxx",
			language.NewFakeLanguage("").WithAllTestFiles(
				func() ([]string, error) {
					return []string{}, nil
				}),
			[]model.CheckPoint{
				model.WarningCheckPoint("no matching test file found"),
			},
		},
		{
			"1 match", "xxx",
			language.NewFakeLanguage("").WithAllTestFiles(
				func() ([]string, error) {
					return []string{"test-file1"}, nil
				}),
			[]model.CheckPoint{
				model.OkCheckPoint("matching test files found:"),
				model.OkCheckPoint("- test-file1"),
			},
		},
		{
			"2 matches", "xxx",
			language.NewFakeLanguage("").WithAllTestFiles(
				func() ([]string, error) {
					return []string{"test-file1", "test-file2"}, nil
				}),
			[]model.CheckPoint{
				model.OkCheckPoint("matching test files found:"),
				model.OkCheckPoint("- test-file1"),
				model.OkCheckPoint("- test-file2"),
			},
		},
		{
			"matching error", "xxx",
			language.NewFakeLanguage("").WithAllTestFiles(
				func() ([]string, error) {
					return nil, errors.New("some error")
				}),
			[]model.CheckPoint{
				model.ErrorCheckPoint("some error"),
			},
		},
		{
			"unreachable directory error and no matching file", "xxx",
			language.NewFakeLanguage("").WithAllTestFiles(
				func() ([]string, error) {
					err := language.UnreachableDirectoryError{}
					err.Add("dir1")
					return nil, &err
				}),
			[]model.CheckPoint{
				model.WarningCheckPoint("cannot access test directory dir1"),
				model.WarningCheckPoint("no matching test file found"),
			},
		},
		{
			"unreachable directory error and one matching file", "xxx",
			language.NewFakeLanguage("").WithAllTestFiles(
				func() ([]string, error) {
					err := language.UnreachableDirectoryError{}
					err.Add("dir1")
					return []string{"test-file1"}, &err
				}),
			[]model.CheckPoint{
				model.WarningCheckPoint("cannot access test directory dir1"),
				model.OkCheckPoint("matching test files found:"),
				model.OkCheckPoint("- test-file1"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithLanguage(test.langName))
			checkEnv.lang = test.lang
			assert.Equal(t, test.expected, checkLanguageTestFiles(p))
		})
	}
}
