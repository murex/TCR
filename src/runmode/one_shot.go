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

package runmode

// OneShot is a type of run mode allowing to run a single TCR cycle and exit
type OneShot struct {
}

// Name returns the name of this run mode
func (OneShot) Name() string {
	return "one-shot"
}

// AutoPushDefault returns the default value of VCS auto-push option with this run mode
func (OneShot) AutoPushDefault() bool {
	return false
}

// IsMultiRole indicates if this run mode supports multiple roles
func (OneShot) IsMultiRole() bool {
	return false
}

// IsInteractive indicates if this run mode allows user interaction
func (OneShot) IsInteractive() bool {
	return false
}

// IsActive indicates if this run mode is actively running TCR
func (OneShot) IsActive() bool {
	return true
}
