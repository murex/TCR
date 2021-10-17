package notification

import (
	"testing"
	"time"
)

func Test_beeep_beep(t *testing.T) {
	tryBeeepBeep(440, 1)
	time.Sleep(2 * time.Second)
}

func Test_beeep_notify(t *testing.T) {
	tryBeeepNotify()
}

func Test_beeep_alert(t *testing.T) {
	tryBeeepAlert()
}
