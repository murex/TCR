/*
Copyright (c) 2024 Murex

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

package commit

import (
	"github.com/murex/tcr/events"
	"os"
	"regexp"
	"strings"
)

const envVariableName = "TCR_COMMIT_MESSAGE_COMMAND"

// NewMessageBuilder returns an instance of commit message builder.
// The returned instance is currently based on TCR_COMMIT_MESSAGE_COMMAND environment variable
// - If set, the contents of this variable is used to create an ExternalMessageBuilder instance,
// which path and arguments are extracted from the env variable value.
// - If not set, a SimpleMessageBuilder instance is returned.
func NewMessageBuilder(header string, event *events.TCREvent, suffix string) MessageBuilder {
	command, useExternalBuilder := os.LookupEnv(envVariableName)

	if useExternalBuilder {
		path, args := extractPathAndArgs(command)
		return NewExternalMessageBuilder(path, args...)
	}
	// Default commit message builder
	return NewSimpleMessageBuilder(header, event, suffix)
}

func extractPathAndArgs(command string) (path string, args []string) {
	// Important: we are assuming here that all arguments are space-free
	// ex: "echo 'hello world'" will produce 2 arguments "'hello" and "world'"
	re := regexp.MustCompile(`[ \t\r\n]+`)
	words := re.Split(strings.TrimSpace(command), -1)
	return words[0], words[1:]
}
