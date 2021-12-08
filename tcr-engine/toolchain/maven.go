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

package toolchain

// MavenToolchain is the toolchain implementation for Maven
type MavenToolchain struct {
}

func (tchn MavenToolchain) reset() {
	//TODO implement me
}

// Name provides the name of the toolchain
func (tchn MavenToolchain) Name() string {
	return "maven"
}

// RunBuild runs the build with this toolchain
func (tchn MavenToolchain) RunBuild() error {
	return runBuild(tchn)
}

// RunTests runs the tests with this toolchain
func (tchn MavenToolchain) RunTests() error {
	return runTests(tchn)
}

// BuildCommandName returns the build command name for this toolchain
func (tchn MavenToolchain) BuildCommandName() string {
	return "mvnw"
}

// BuildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn MavenToolchain) BuildCommandArgs() []string {
	return []string{"test-compile"}
}

// TestCommandName returns the test command name for this toolchain
func (tchn MavenToolchain) TestCommandName() string {
	return "mvnw"
}

// TestCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn MavenToolchain) TestCommandArgs() []string {
	return []string{"test"}
}
