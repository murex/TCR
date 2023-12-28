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

package settings

import (
	"fmt"
	"time"
)

// Below variables are set at build time through -ldflags
var (
	BuildVersion = "v0.0.0-dev"
	BuildOs      = "unknown"
	BuildArch    = "unknown"
	BuildCommit  = "none"
	BuildDate    = time.Time{}.Format(time.RFC3339)
	BuildAuthor  = "unknown"
)

// BuildInfo contains build information in a Label/Value format
type BuildInfo struct {
	Label string
	Value string
}

// GetBuildInfo returns a table with TCR build information
func GetBuildInfo() []BuildInfo {
	var t []BuildInfo
	t = append(t, BuildInfo{"Version", BuildVersion})
	t = append(t, BuildInfo{"OS Family", BuildOs})
	t = append(t, BuildInfo{"Architecture", BuildArch})
	t = append(t, BuildInfo{"Commit", BuildCommit})
	t = append(t, BuildInfo{"Build Date", BuildDate})
	t = append(t, BuildInfo{"Built By", BuildAuthor})
	return t
}

// PrintBuildInfo prints out application's build information and exits
func PrintBuildInfo() {
	for _, buildInfo := range GetBuildInfo() {
		fmt.Printf("- %s:\t%s\n", buildInfo.Label, buildInfo.Value)
	}
}
