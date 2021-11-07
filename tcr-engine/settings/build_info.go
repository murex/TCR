/*
Copyright (c) 2021 Murex

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

import "fmt"

// Below variables are set at build time through -ldflags
var (
	BuildVersion = "dev"
	BuildOs      = "unknown"
	BuildArch    = "unknown"
	BuildCommit  = "none"
	BuildDate    = "unknown"
	BuildAuthor  = "unknown"
)

// GetBuildInfo returns TCR build information as a map
func GetBuildInfo() map[string]string {
	var m = make(map[string]string)
	m["Version"] = BuildVersion
	m["OS Family"] = BuildOs
	m["Architecture"] = BuildArch
	m["Commit"] = BuildCommit
	m["Build Date"] = BuildDate
	m["Built By"] = BuildAuthor
	return m
}

// PrintBuildInfo prints information related to the build
func PrintBuildInfo() {
	for key, value := range GetBuildInfo() {
		fmt.Printf("- %s:\t%s\n", key, value)
	}
}
