package notification

import "github.com/gen2brain/beeep"

func tryBeeepBeep(frequency float64, duration int) {
	err := beeep.Beep(frequency, duration)
	if err != nil {
		panic(err)
	}
}

func tryBeeepNotify() {
	err := beeep.Notify("TCR Mob Timer", "Still 3m left", "")
	if err != nil {
		panic(err)
	}
}

func tryBeeepAlert() {
	err := beeep.Alert("TCR Mob Timer", "Time to rotate!", "assets/warning.png")
	if err != nil {
		panic(err)
	}
}
