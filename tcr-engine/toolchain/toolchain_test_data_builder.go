//go:build test_helper

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

package toolchain

// AToolchain is a test data builder for type Toolchain
func AToolchain(toolchainBuilders ...func(tchn *Toolchain)) *Toolchain {
	tchn := New("default-toolchain", []Command{*ACommand()}, []Command{*ACommand()})

	for _, build := range toolchainBuilders {
		build(tchn)
	}
	return tchn
}

// WithName sets the name of the created toolchain to name
func WithName(name string) func(tchn *Toolchain) {
	return func(tchn *Toolchain) { tchn.name = name }
}

// WithNoBuildCommand creates a toolchain with no build command defined
func WithNoBuildCommand() func(tchn *Toolchain) {
	return func(tchn *Toolchain) { tchn.buildCommands = nil }
}

// WithBuildCommand adds the provided command as a build command
func WithBuildCommand(command *Command) func(tchn *Toolchain) {
	return func(tchn *Toolchain) {
		tchn.buildCommands = append(tchn.buildCommands, *command)
	}
}

// WithNoTestCommand creates a toolchain with no test command defined
func WithNoTestCommand() func(tchn *Toolchain) {
	return func(tchn *Toolchain) { tchn.testCommands = nil }
}

// WithTestCommand adds the provided command as a test command
func WithTestCommand(command *Command) func(tchn *Toolchain) {
	return func(tchn *Toolchain) {
		tchn.testCommands = append(tchn.testCommands, *command)
	}
}
