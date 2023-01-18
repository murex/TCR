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

package p4

import (
	"bufio"
	"bytes"
	"github.com/murex/tcr/vcs/shell"
	"strings"
)

func init() {
	shell.NewCommandFunc = shell.NewCommand
}

func newP4Command(params ...string) shell.Command {
	return shell.NewCommandFunc("p4", params...)
}

func newP4CommandImpl(params ...string) *shell.CommandImpl {
	return shell.NewCommandImpl("p4", params...)
}

// IsP4CommandAvailable indicates if p4 command is available on local machine
func IsP4CommandAvailable() bool {
	return newP4Command().IsInPath()
}

// GetP4CommandPath returns the path to p4 command on this machine
func GetP4CommandPath() string {
	return newP4Command().GetFullPath()
}

// GetP4CommandVersion returns the version of p4 command on this machine
func GetP4CommandVersion() string {
	p4Output, err := runP4Command("-V")
	if err != nil {
		return "unknown"
	}
	scanner := bufio.NewScanner(bytes.NewReader(p4Output))
	for scanner.Scan() {
		if strings.Index(scanner.Text(), "Rev.") == 0 {
			return strings.Split(scanner.Text(), " ")[1]
		}
	}
	return ""
}

// GetP4UserName returns the user name retrieved from local p4 configuration
func GetP4UserName() string {
	return getP4ConfigValue("P4USER")
}

// GetP4ClientName returns the client name retrieved from the local p4 configuration
func GetP4ClientName() string {
	return getP4ConfigValue("P4CLIENT")
}

// GetRootDir retrieves the local root directory for the depot's workspace
func GetRootDir() (string, error) {
	root, err := runP4Command("-F", "%clientRoot%", "-ztag", "info")
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(root), "\r\n"), nil
}

func getP4ConfigValue(variable string) string {
	p4Output, err := runP4Command("set", "-q", variable)

	if err != nil || p4Output == nil || len(bytes.Trim(p4Output, "\r\n")) == 0 {
		return "not set"
	}
	scanner := bufio.NewScanner(bytes.NewReader(p4Output))
	scanner.Scan()
	return strings.TrimPrefix(scanner.Text(), variable+"=")
}

// traceP4Command calls p4 command and reports its output traces
func traceP4Command(params ...string) error {
	return newP4Command().Trace(params...)
}

// runP4Command calls p4 command in a separate process and returns its output traces
func runP4Command(params ...string) (output []byte, err error) {
	return newP4Command().Run(params...)
}

// tracePipedP4Command calls p4 command, pipes it to pipedTo command, and reports its output traces
func tracePipedP4Command(pipedTo shell.Command, params ...string) error {
	return newP4Command().TraceAndPipe(pipedTo, params...)
}

// runPipedP4Command calls p4 command, pipes it to pipedTo command, and reports its output traces
func runPipedP4Command(pipedTo shell.Command, params ...string) (output []byte, err error) {
	return newP4Command().RunAndPipe(pipedTo, params...)
}
