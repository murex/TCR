package timer

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var testTimeout time.Duration
var testTickPeriod time.Duration

func TestMain(m *testing.M) {
	if !flag.Parsed() {
		flag.Parse()
	}
	// Most tests in this file are designed so that reminders fire twice, and then go in time out.
	// We differentiate CI and local machine to optimize test speed execution when run on local machine
	// while not failing when run on CI (which runs slower)
	if testing.Short() {
		fmt.Println("Running tests with fast timers")
		testTimeout = 100 * time.Millisecond
		testTickPeriod = 40 * time.Millisecond
	} else {
		fmt.Println("Running tests with long timers")
		testTimeout = 1000 * time.Millisecond
		testTickPeriod = 400 * time.Millisecond
	}
	// Run tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

// Timeout

func Test_default_timeout_is_5_min(t *testing.T) {
	r := New(0, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	assert.Equal(t, 5*time.Minute, r.timeout)
}

func Test_init_with_non_default_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	assert.Equal(t, testTimeout, r.timeout)
}

func Test_ticking_stops_after_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTimeout * 2)
	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// Tick Period

func Test_default_tick_period_is_1_min(t *testing.T) {
	r := New(testTimeout, 0, func(tc TickContext) {}, func(tc TickContext) {})
	assert.Equal(t, 1*time.Minute, r.tickPeriod)
}

func Test_init_with_non_default_tick_period(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	assert.Equal(t, testTickPeriod, r.tickPeriod)
}

// Starting PeriodicReminder

func Test_start_reminder(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	time.Sleep(testTimeout)
	assert.Equal(t, 0, r.tickCounter)
	r.Start()
	time.Sleep(testTimeout)
	assert.Equal(t, 2, r.tickCounter)
}

// Stopping PeriodicReminder

func Test_stop_reminder_before_1st_tick(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTickPeriod / 2)
	r.Stop()
	time.Sleep(testTimeout)

	assert.Equal(t, 0, r.tickCounter)
	assert.Equal(t, StoppedAfterInterruption, r.state)
}

func Test_stop_reminder_between_1st_and_2nd_tick(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTickPeriod + testTickPeriod/2)
	r.Stop()
	time.Sleep(testTimeout)

	assert.Equal(t, 1, r.tickCounter)
	assert.Equal(t, StoppedAfterInterruption, r.state)
}

func Test_stop_reminder_after_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTimeout * 2)
	r.Stop()

	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// PeriodicReminder tick counter

func Test_can_track_number_of_ticks_fired(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	r.Start()
	assert.Equal(t, 0, r.tickCounter)
	time.Sleep(testTickPeriod / 2)
	assert.Equal(t, 0, r.tickCounter)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 1, r.tickCounter)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 2, r.tickCounter)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 2, r.tickCounter)
	assert.Equal(t, StoppedAfterTimeOut, r.state)
}

// PeriodicReminder callback function

func Test_callback_function_can_know_current_tick_index(t *testing.T) {
	var index int
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {
		index = tc.index
	}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTickPeriod + testTickPeriod/2)
	assert.Equal(t, 0, index)
	time.Sleep(testTickPeriod)
	assert.Equal(t, 1, index)
}

func Test_callback_function_can_know_timestamp(t *testing.T) {
	var ts [2]time.Time
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {
		ts[tc.index] = tc.timestamp
	}, func(tc TickContext) {})
	tsStart := time.Now()
	r.Start()
	time.Sleep(testTimeout)
	tsEnd := time.Now()

	assert.True(t, tsStart.Before(ts[0]))
	assert.True(t, ts[0].Before(ts[1]))
	assert.True(t, ts[1].Before(tsEnd))
}

func Test_callback_function_can_know_elapsed_time_since_start(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {
		var expected = testTickPeriod * time.Duration(tc.index+1)
		assert.Equal(t, expected, tc.elapsed)
	}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTimeout)
}

func Test_callback_function_can_know_remaining_time_until_end(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {
		var expected = testTimeout - testTickPeriod*time.Duration(tc.index+1)
		assert.Equal(t, expected, tc.remaining)
	}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTimeout)
}

func Test_callback_function_can_know_max_index_value(t *testing.T) {
	var expected = int(testTimeout/testTickPeriod) - 1
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {
		assert.Equal(t, expected, tc.indexMax)
	}, func(tc TickContext) {})
	r.Start()
	time.Sleep(testTimeout)
}

// Time elapsed since timer started

func Test_retrieving_time_elapsed_since_timer_started(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	// Before calling Start(), time elapsed should stick to 0
	assert.Zero(t, r.GetElapsedTime())
	r.Start()
	time.Sleep(testTimeout / 2)
	// While timer is running, total time elapsed is time spent since Start()
	assert.InEpsilon(t, testTimeout/2, r.GetElapsedTime(), 0.3)
	time.Sleep(testTimeout)
	// When timer is done, time elapsed should stopTicking increasing
	assert.InEpsilon(t, testTimeout, r.GetElapsedTime(), 0.3)
}

// Time remaining until timer ends

func Test_retrieving_time_remaining_until_timer_ends(t *testing.T) {
	r := New(testTimeout, testTickPeriod, func(tc TickContext) {}, func(tc TickContext) {})
	// Before calling Start(), time remaining should stick to timeout
	assert.Equal(t, testTimeout, r.GetRemainingTime())
	r.Start()
	time.Sleep(testTimeout / 2)
	// While timer is running, total time remaining is timeout - time spent since Start()
	assert.InEpsilon(t, testTimeout/2, r.GetRemainingTime(), 0.3)
	time.Sleep(testTimeout)
	// When timer is done, time remaining should be 0
	assert.Zero(t, r.GetRemainingTime())
}
