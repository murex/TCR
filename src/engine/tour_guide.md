# Driver Round

## Let's start as a driver

```go
// RunAsDriver tells TCR engine to start running with driver role
func (tcr *TCREngine) RunAsDriver() {
	// Force previous role to quit if needed
	if tcr.GetCurrentRole() != nil {
		tcr.Stop()
	}
	// Prepare the timer
```
[Link to code](https://github.com/murex/TCR/blob/main/src/engine/././tcr.go#L403-L410)

## We made some code changes, let's see what happens ⏳

```go
// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (tcr *TCREngine) RunTCRCycle() {
	status.RecordState(status.Ok)
	if tcr.build().Failed() {
		return
	}
```
[Link to code](https://github.com/murex/TCR/blob/main/src/engine/././tcr.go#L517-L524)

## The tests are running, fingers crossed! 🤞

```go
func (tcr *TCREngine) test() (result toolchain.TestCommandResult) {
	report.PostInfo("Running Tests")
	result = tcr.toolchain.RunTests()
	if result.Failed() {
		status.RecordState(status.TestFailed)
		report.PostErrorWithEmphasis(testFailureMessage)
	} else {
```
[Link to code](https://github.com/murex/TCR/blob/main/src/engine/././tcr.go#L570-L577)

## Damned, the tests failed, the code is reverted! 😭

```go
func (tcr *TCREngine) revertSrcFiles() {
	diffs, err := tcr.vcs.Diff()
	tcr.handleError(err, false, status.VCSError)
	if err != nil {
		return
	}
	var reverted int
```
[Link to code](https://github.com/murex/TCR/blob/main/src/engine/././tcr.go#L647-L654)

## Let's try again, hopefully it will work this time ⏳

```go
// RunTCRCycle is the core of TCR engine: e.g. it runs one test && commit || revert cycle
func (tcr *TCREngine) RunTCRCycle() {
	status.RecordState(status.Ok)
	if tcr.build().Failed() {
		return
	}
	result := tcr.test()
```
[Link to code](https://github.com/murex/TCR/blob/main/src/engine/././tcr.go#L518-L525)

## Great! The tests are passing this time. The changes are committed 😀

```go
func (tcr *TCREngine) commit(event events.TCREvent) {
	report.PostInfo("Committing changes on ", tcr.vcs.SessionSummary())
	var err error
	err = tcr.vcs.Add()
	tcr.handleError(err, false, status.VCSError)
	if err != nil {
		return
```
[Link to code](https://github.com/murex/TCR/blob/main/src/engine/././tcr.go#L583-L590)

