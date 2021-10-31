/*
Copyright (c) 2021 Murex

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

package settings

import "time"

// Default values

const (
	// ApplicationName is the name of the application
	ApplicationName = "TCR"
	// DefaultPollingPeriod is the waiting time between 2 consecutive calls to git pull when running as Navigator
	DefaultPollingPeriod = 2 * time.Second
	// DefaultMobTurnDuration is the default duration for a mob turn
	DefaultMobTurnDuration = 5 * time.Minute
	// DefaultInactivityPeriod is the default inactivity period until TCR sends an inactivity teaser message
	DefaultInactivityPeriod = 1 * time.Minute
	// DefaultInactivityTimeout is the default timeout after which TCR stops sending inactivity teaser messages
	DefaultInactivityTimeout = 5 * time.Minute
)
