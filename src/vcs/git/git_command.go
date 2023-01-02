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

package git

import (
	"bufio"
	"bytes"
	"github.com/murex/tcr/vcs/cmd"
	"strings"
)

func newGitCommand() *cmd.ShellCommand {
	return cmd.New("git")
}

// IsGitCommandAvailable indicates if git command is available on local machine
func IsGitCommandAvailable() bool {
	return newGitCommand().IsInPath()
}

// GetGitCommandPath returns the path to git command on this machine
func GetGitCommandPath() string {
	return newGitCommand().GetFullPath()
}

// GetGitCommandVersion returns the version of git command on this machine
func GetGitCommandVersion() string {
	gitOutput, err := runGitCommand("version")
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
	gitOutput, err := runGitCommand("config", variable)
	if err != nil || gitOutput == nil || len(gitOutput) == 0 {
		return "not set"
	}
	scanner := bufio.NewScanner(bytes.NewReader(gitOutput))
	scanner.Scan()
	return scanner.Text()
}

// traceGitCommand calls git command and reports its output traces
func traceGitCommand(params ...string) error {
	return newGitCommand().Trace(params...)
}

// runGitCommand calls git command in a separate process and returns its output traces
func runGitCommand(params ...string) (output []byte, err error) {
	return newGitCommand().Run(params...)
}
