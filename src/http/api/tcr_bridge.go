package api

import "github.com/murex/tcr/engine"

var tcr engine.TCRInterface

func SetTCRInstance(instance engine.TCRInterface) {
	tcr = instance
}
