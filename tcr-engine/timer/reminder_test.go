package timer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testTimeout = 250 * time.Millisecond
const testTickPeriod = 100 * time.Millisecond

// Timeout

func Test_default_timeout_is_5_min(t *testing.T) {
	r := New(0, testTickPeriod, nil)
	assert.Equal(t, 5 * time.Minute, r.timeout)
}

func Test_init_with_non_default_timeout(t *testing.T) {
	r := New(testTimeout, testTickPeriod, nil)
	assert.Equal(t, testTimeout, r.timeout)
}

// Tick Period

func Test_default_tick_period_is_1_min(t *testing.T) {
	r := New(testTimeout, 0, nil)
	assert.Equal(t, 1 * time.Minute, r.tickPeriod)
}

func Test_init_with_non_default_tick_period(t *testing.T) {
	r := New(testTimeout, testTickPeriod, nil)
	assert.Equal(t, testTickPeriod, r.tickPeriod)
}

// Start Reminder

func Test_start_reminder(t *testing.T) {
	var nbTicks = 0
	r := New(testTimeout, testTickPeriod, func(t time.Time) {
		nbTicks++
	})
	time.Sleep(testTimeout)
	assert.Equal(t, 0, nbTicks)
	r.Start()
	time.Sleep(testTimeout)
	assert.Equal(t, 2, nbTicks)
}

// Stop Reminder

func Test_stop_reminder_before_1st_tick(t *testing.T) {
	var nbTicks = 0
	r := New(testTimeout, testTickPeriod, func(t time.Time) {
		nbTicks++
	})
	fmt.Println("A ", time.Now())
	r.Start()
	fmt.Println("B ", time.Now())
	time.Sleep(testTickPeriod / 2)
	fmt.Println("C ", time.Now())
	r.Stop()
	fmt.Println("D ", time.Now())
	time.Sleep(testTimeout)
	fmt.Println("E ", time.Now())
	assert.Equal(t, 0, nbTicks)
}

func Test_stop_reminder_between_1st_and_2nd_tick(t *testing.T) {
	var nbTicks = 0
	r := New(testTimeout, testTickPeriod, func(t time.Time) {
		nbTicks++
	})
	fmt.Println("A ", time.Now())
	r.Start()
	fmt.Println("B ", time.Now())
	time.Sleep(testTickPeriod + testTickPeriod / 2)
	fmt.Println("C ", time.Now())
	r.Stop()
	fmt.Println("D ", time.Now())
	time.Sleep(testTimeout)
	fmt.Println("E ", time.Now())
	assert.Equal(t, 1, nbTicks)
}

//func Test_stop_reminder_after_timeout(t *testing.T) {
//	var nbTicks = 0
//	r := New(testTimeout, testTickPeriod, func(t time.Time) {
//		nbTicks++
//	})
//	fmt.Println("A ", time.Now())
//	r.Start()
//	fmt.Println("B ", time.Now())
//	time.Sleep(testTimeout * 2)
//	fmt.Println("C ", time.Now())
//	r.Stop()
//	fmt.Println("D ", time.Now())
//	time.Sleep(testTimeout)
//	fmt.Println("E ", time.Now())
//	assert.Equal(t, 2, nbTicks)
//}