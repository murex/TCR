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

package desktop

import (
	"fmt"
	"net/url"

	"github.com/murex/tcr/report"
	"github.com/pkg/browser"
)

const hostname = "127.0.0.1"

// OpenBrowser opens a browser on localhost on the provided port number
func OpenBrowser(portNumber int) {
	u := browserURL(portNumber)
	err := browser.OpenURL(u)
	if err != nil {
		report.PostWarning("Could not open ", u, ": ", err.Error())
	}
}

func browserURL(portNumber int) string {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", hostname, portNumber),
		Path:   "/",
	}
	return u.String()
}
