package engine

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/murex/tcr/params"
	"github.com/murex/tcr/report"
	"github.com/murex/tcr/runmode"
	"testing"
	"time"
)

// tcrCtxKey is the key used to store the available godogs in the context.Context.
type tcrCtxKey struct{}
type snifferCtxKey struct{}

const TIMER_DURATION = 1 * time.Minute

func beforeScenario(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	sniffer := report.NewSniffer(
		func(msg report.Message) bool {
			return msg.Type.Category == report.Info &&
				msg.Payload.ToString() == "Timer duration is "+TIMER_DURATION.String()
		},
	)
	sniffer.Start()
	return context.WithValue(ctx, snifferCtxKey{}, sniffer), nil
}

func aTCREngineWithAMobTimerEnabled(ctx context.Context) (context.Context, error) {
	var tcr *TCREngine
	tcr, _ = initTCREngineWithFakes(params.AParamSet(
		params.WithRunMode(runmode.Mob{}),
		params.WithMobTimerDuration(TIMER_DURATION),
	), nil, nil, nil)

	return context.WithValue(ctx, tcrCtxKey{}, tcr), nil
}
func theTCREngineStarts(ctx context.Context) error {
	tcr := ctx.Value(tcrCtxKey{}).(*TCREngine)

	tcr.RunAsDriver()
	time.Sleep(1 * time.Millisecond)

	return nil
}

func theTimerDurationShouldBeTracedExactlyOnce(ctx context.Context) error {
	tcr := ctx.Value(tcrCtxKey{}).(*TCREngine)
	sniffer := ctx.Value(snifferCtxKey{}).(*report.Sniffer)

	tcr.Stop()

	sniffer.Stop()
	timerDurationTracesCount := sniffer.GetMatchCount()

	if timerDurationTracesCount != 1 {
		return fmt.Errorf("expected exactly one timer duration message, but got %d", timerDurationTracesCount)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(beforeScenario)

	ctx.Step(`^a TCR engine with a mob timer enabled$`, aTCREngineWithAMobTimerEnabled)
	ctx.Step(`^the TCR engine starts$`, theTCREngineStarts)
	ctx.Step(`^the timer duration should be traced exactly once$`, theTimerDurationShouldBeTracedExactlyOnce)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
