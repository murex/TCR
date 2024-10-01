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

package variant

import (
	"fmt"
	"strings"
)

// UnsupportedVariantError is returned when the provided Variant name is not supported.
type UnsupportedVariantError struct {
	variantName string
}

// Error returns the error description
func (e *UnsupportedVariantError) Error() string {
	return fmt.Sprintf("variant not supported: \"%s\"", e.variantName)
}

// Variant represents the possible values for the TCR Variant.
// These values are inspired by the following blog-post:
// https://medium.com/@tdeniffel/tcr-variants-test-commit-revert-bf6bd84b17d3
type Variant string

// Recognized variant values
const (
	Relaxed       Variant = "relaxed"
	BTCR          Variant = "btcr"
	Introspective Variant = "introspective"
)

var recognized = []Variant{Relaxed, BTCR, Introspective}

// Select returns a variant instance for the provided name.
// It returns an UnsupportedVariantError if the name is not recognized as a
// valid variant name.
func Select(name string) (*Variant, error) {
	for _, variant := range recognized {
		if strings.EqualFold(name, variant.Name()) {
			return &variant, nil
		}
	}
	return nil, &UnsupportedVariantError{name}
}

// Name returns the variant name
func (v Variant) Name() string {
	return string(v)
}
