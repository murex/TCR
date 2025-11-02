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

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_building_commit_messages(t *testing.T) {
	tests := []struct {
		desc          string
		commitMessage CommitMessage
		withEmoji     bool
		expected      string
	}{
		{
			"Passing Message",
			CommitMessage{
				Emoji:       '✅',
				Tag:         "[TCR - PASSED]",
				Description: "tests passing",
			},
			true,
			passedCommitMessage,
		}, {
			"Failing Message",
			CommitMessage{
				Emoji:       '❌',
				Tag:         "[TCR - FAILED]",
				Description: "tests failing",
			},
			true,
			failedCommitMessage,
		}, {
			"Reverting Message",
			CommitMessage{
				Emoji:       '⏪',
				Tag:         "[TCR - REVERTED]",
				Description: "revert changes",
			},
			true,
			revertedCommitMessage,
		}, {
			"Passing Message without Emoji",
			CommitMessage{
				Emoji:       '✅',
				Tag:         "[TCR - PASSED]",
				Description: "tests passing",
			},
			false,
			"[TCR - PASSED] tests passing",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, test.commitMessage.toString(test.withEmoji))
		})
	}
}
