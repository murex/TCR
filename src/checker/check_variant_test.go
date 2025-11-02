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

package checker

import (
	"testing"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/stretchr/testify/assert"
)

func Test_check_variant_configuration(t *testing.T) {
	assertCheckGroupRunner(t,
		checkVariantConfiguration,
		&checkVariantRunners,
		*params.AParamSet(),
		"TCR variant configuration")
}

func Test_check_variant_selection(t *testing.T) {
	tests := []struct {
		desc        string
		variantName string
		expected    []model.CheckPoint
	}{
		{
			"empty", "",
			[]model.CheckPoint{
				model.ErrorCheckPoint("no variant is selected"),
			},
		},
		{
			"unknown", "unknown-variant",
			[]model.CheckPoint{
				model.ErrorCheckPoint("selected variant is not supported: \"unknown-variant\""),
			},
		},
		{
			"Original", "original",
			[]model.CheckPoint{
				model.ErrorCheckPoint("original variant is not yet supported"),
			},
		},
		{
			"Introspective", "introspective",
			[]model.CheckPoint{
				model.OkCheckPoint("selected variant is introspective"),
			},
		},
		{
			"BTCR", "btcr",
			[]model.CheckPoint{
				model.OkCheckPoint("selected variant is btcr"),
			},
		},
		{
			"Relaxed", "relaxed",
			[]model.CheckPoint{
				model.OkCheckPoint("selected variant is relaxed"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := *params.AParamSet(params.WithVariant(test.variantName))
			assert.Equal(t, test.expected, checkVariantSelection(p))
		})
	}
}
