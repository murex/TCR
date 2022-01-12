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

// ACommand is a test data builder for type Command
func ACommand(commandBuilders ...func(command *Command)) *Command {
	command := &Command{
		Os:        GetAllOsNames(),
		Arch:      GetAllArchNames(),
		Path:      "",
		Arguments: []string{},
	}

	for _, build := range commandBuilders {
		build(command)
	}
	return command
}

// WithPath sets the command path for the created command
func WithPath(path string) func(command *Command) {
	return func(command *Command) { command.Path = path }
}

// WithOs adds the provided OS to the list of supported OS's for this command
func WithOs(os OsName) func(command *Command) {
	return func(command *Command) { command.Os = append(command.Os, os) }
}

// WithNoOs creates a command with no supported OS
func WithNoOs() func(command *Command) {
	return func(command *Command) { command.Os = nil }
}

// WithArch adds the provided OS architecture to the list of supported architectures for this command
func WithArch(arch ArchName) func(command *Command) {
	return func(command *Command) { command.Arch = append(command.Arch, arch) }
}

// WithNoArch creates a command with no supported OS architecture
func WithNoArch() func(command *Command) {
	return func(command *Command) { command.Arch = nil }
}
