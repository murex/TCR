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

// CmakeToolchain is the toolchain implementation for CMake
type CmakeToolchain struct{}

func (tchn CmakeToolchain) reset() {
	//TODO implement me
}

// Name provides the name of the toolchain
func (tchn CmakeToolchain) Name() string {
	return "cmake"
}

// RunBuild runs the build with this toolchain
func (tchn CmakeToolchain) RunBuild() error {
	return runBuild(tchn)
}

// RunTests runs the tests with this toolchain
func (tchn CmakeToolchain) RunTests() error {
	return runTests(tchn)
}

// BuildCommandArgs returns a table with the list of build command arguments for this toolchain
func (tchn CmakeToolchain) BuildCommandArgs() []string {
	return []string{"--build", "build", "--config", "Debug"}
}

// TestCommandArgs returns a table with the list of test command arguments for this toolchain
func (tchn CmakeToolchain) TestCommandArgs() []string {
	// Important: This (--test-dir option) requires using cmake 3.20 version or higher
	return []string{"--output-on-failure", "--test-dir", "build", "--build-config", "Debug"}
}
