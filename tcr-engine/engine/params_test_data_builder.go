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

package engine

import (
	"github.com/murex/tcr/tcr-engine/runmode"
	"time"
)

// AParamSet is a test data builder for type Params
func AParamSet(builders ...func(params *Params)) *Params {
	params := &Params{
		ConfigDir:       "",
		BaseDir:         "",
		Language:        "",
		Toolchain:       "",
		MobTurnDuration: 0,
		AutoPush:        false,
		PollingPeriod:   0,
		Mode:            runmode.Check{},
	}

	for _, build := range builders {
		build(params)
	}
	return params
}

// WithBaseDir sets the provided dirName as the base directory
func WithBaseDir(dirName string) func(params *Params) {
	return func(params *Params) {
		params.BaseDir = dirName
	}
}

// WithConfigDir sets the provided dirName as the configuration directory
func WithConfigDir(dirName string) func(params *Params) {
	return func(params *Params) {
		params.ConfigDir = dirName
	}
}

// WithLanguage sets the provided value as the language name
func WithLanguage(name string) func(params *Params) {
	return func(params *Params) {
		params.Language = name
	}
}

// WithToolchain sets the provided value as the toolchain name
func WithToolchain(name string) func(params *Params) {
	return func(params *Params) {
		params.Toolchain = name
	}
}

// WithPollingPeriod sets the provided value as the git polling period
func WithPollingPeriod(period time.Duration) func(params *Params) {
	return func(params *Params) {
		params.PollingPeriod = period
	}
}

// WithMobTimerDuration sets the provided value as the mob timer duration
func WithMobTimerDuration(duration time.Duration) func(params *Params) {
	return func(params *Params) {
		params.MobTurnDuration = duration
	}
}

// WithAutoPush sets auto-push flag to true
func WithAutoPush() func(params *Params) {
	return func(params *Params) {
		params.AutoPush = true
	}
}

// WithNoAutoPush sets auto-push flag to false
func WithNoAutoPush() func(params *Params) {
	return func(params *Params) {
		params.AutoPush = false
	}
}
