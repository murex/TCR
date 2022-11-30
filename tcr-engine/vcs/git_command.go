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

package vcs

import (
	"bufio"
	"bytes"
	"github.com/codeskyblue/go-sh"
	"github.com/murex/tcr/report"
	"strings"
)

// IsGitCommandAvailable indicates if git command is available on local machine
func IsGitCommandAvailable() bool {
	return isCommandAvailable("git")
}

// GetGitCommandPath returns the path to git command on this machine
func GetGitCommandPath() string {
	return getCommandPath("git")
}

// GetGitCommandVersion returns the version of git command on this machine
func GetGitCommandVersion() string {
	gitOutput, err := runGitCommand([]string{"version"})
	if err != nil {
		return "unknown"
	}
	scanner := bufio.NewScanner(bytes.NewReader(gitOutput))
	scanner.Scan()
	return strings.SplitAfter(scanner.Text(), "git version ")[1]
}

// GetGitUserName returns the user name retrieved from local git configuration
func GetGitUserName() string {
	return getGitConfigValue("user.name")
}

func getGitConfigValue(variable string) string {
	gitOutput, err := runGitCommand([]string{"config", variable})
	if err != nil || gitOutput == nil || len(gitOutput) == 0 {
		return "not set"
	}
	scanner := bufio.NewScanner(bytes.NewReader(gitOutput))
	scanner.Scan()
	return scanner.Text()
}

// traceGitCommand calls git command and reports its output traces
func traceGitCommand(params []string) error {
	output, err := runGitCommand(params)
	if len(output) > 0 {
		report.PostText(string(output))
	}
	return err
}

// runGitCommand calls git command in a separate process and returns its output traces
func runGitCommand(params []string) (output []byte, err error) {
	return sh.Command("git", params).CombinedOutput()
}
