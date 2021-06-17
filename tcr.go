package main

import (
	"github.com/mengdaming/tcr/trace"
)

func main() {
	trace.Info("Hello, world.")
	trace.HorizontalLine()
	trace.Warning("Hello, world.")
	trace.Error("Hello, world.")
}

// TODO CLI using Cobra: https://dzone.com/articles/how-to-create-a-cli-in-go-in-few-minutes