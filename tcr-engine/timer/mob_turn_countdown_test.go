package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_tick_period_for_timeout_lower_than_10s_is_1s(t *testing.T) {
	expected := 1 * time.Second
	assert.Equal(t, expected, findBestTickPeriodFor(1*time.Second))
	assert.Equal(t, expected, findBestTickPeriodFor(10*time.Second))
}

func Test_tick_period_for_timeout_between_10s_and_1m_is_10s(t *testing.T) {
	expected := 10 * time.Second
	assert.Equal(t, expected, findBestTickPeriodFor(11*time.Second))
	assert.Equal(t, expected, findBestTickPeriodFor(1*time.Minute))
}

func Test_tick_period_for_timeout_between_1m_and_10m_is_1m(t *testing.T) {
	expected := 1 * time.Minute
	assert.Equal(t, expected, findBestTickPeriodFor(1*time.Minute+1*time.Second))
	assert.Equal(t, expected, findBestTickPeriodFor(10*time.Minute))
}
