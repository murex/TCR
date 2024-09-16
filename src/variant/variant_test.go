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

package variant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_get_variant_name(t *testing.T) {
	tests := []struct {
		desc     string
		variant  Variant
		expected string
	}{
		{"relaxed", Relaxed, "relaxed"},
		{"btcr", BTCR, "btcr"},
		{"introspective", Introspective, "introspective"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expected, test.variant.Name())
		})
	}
}

func Test_select_variant(t *testing.T) {
	relaxed, btcr, introspective := Relaxed, BTCR, Introspective
	tests := []struct {
		name            string
		expectedVariant *Variant
		expectedError   error
	}{
		{"relaxed", &relaxed, nil},
		{"btcr", &btcr, nil},
		{"introspective", &introspective, nil},
		{"unknown", nil, &UnsupportedVariantError{"unknown"}},
		{"", nil, &UnsupportedVariantError{""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			variant, err := Select(test.name)
			assert.Equal(t, test.expectedVariant, variant)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_unsupported_variant_message_format(t *testing.T) {
	err := UnsupportedVariantError{"some-variant"}
	assert.Equal(t, "variant not supported: \"some-variant\"", err.Error())
}
