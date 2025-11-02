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
	"strings"

	"github.com/murex/tcr/checker/model"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/variant"
)

var checkVariantRunners []checkPointRunner

func init() {
	checkVariantRunners = []checkPointRunner{
		checkVariantSelection,
	}
}

func checkVariantConfiguration(p params.Params) (cg *model.CheckGroup) {
	cg = model.NewCheckGroup("TCR variant configuration")
	for _, runner := range checkVariantRunners {
		cg.Add(runner(p)...)
	}
	return cg
}

func checkVariantSelection(p params.Params) (cp []model.CheckPoint) {
	switch variantName := strings.ToLower(p.Variant); variantName {
	case variant.Relaxed.Name(), variant.BTCR.Name(), variant.Introspective.Name():
		cp = append(cp, model.OkCheckPoint("selected variant is ", variantName))
	case "original":
		cp = append(cp, model.ErrorCheckPoint("original variant is not yet supported"))
	case "":
		cp = append(cp, model.ErrorCheckPoint("no variant is selected"))
	default:
		cp = append(cp, model.ErrorCheckPoint("selected variant is not supported: \"", variantName, "\""))
	}
	return cp
}
