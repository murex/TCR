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
	"bufio"
	"bytes"
	"github.com/murex/tcr/vcs/shell"
)

// ExternalMessageBuilder provides implementation for an external commit message builder.
// This requires invoking an external command that is in charge of generating
// the commit messages
type ExternalMessageBuilder struct {
	commandPath string
	commandArgs []string
}

// NewExternalMessageBuilder creates a new ExternalMessageBuilder instance
func NewExternalMessageBuilder(commandPath string, commandArgs ...string) *ExternalMessageBuilder {
	return &ExternalMessageBuilder{
		commandPath,
		commandArgs,
	}
}

// GenerateMessage generates a commit message
func (emb *ExternalMessageBuilder) GenerateMessage() (message []string, err error) {
	cmd := shell.NewCommand(emb.commandPath, emb.commandArgs...)
	output, err := cmd.Run()
	if err != nil || output == nil || len(bytes.Trim(output, "\r\n")) == 0 {
		return []string{}, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		message = append(message, scanner.Text())
	}
	return message, nil
}
