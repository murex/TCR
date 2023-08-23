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

package runmode

// Retro is a type of run mode allowing to generate a retrospective template
type Retro struct {
}

// Name returns the name of this run mode
func (Retro) Name() string {
	return "retro"
}

// AutoPushDefault returns the default value of VCS auto-push option with this run mode
func (Retro) AutoPushDefault() bool {
	return false
}

// NeedsCountdownTimer indicates if a countdown timer is needed with this run mode
func (Retro) NeedsCountdownTimer() bool {
	return false
}

// IsInteractive indicates if this run mode allows user interaction
func (Retro) IsInteractive() bool {
	return false
}

// IsActive indicates if this run mode is actively running TCR
func (Retro) IsActive() bool {
	return false
}
