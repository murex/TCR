# Debugging TCR

## Running and debugging From GoLand or IntelliJ

There is run configuration in [.run](../.run) directory to launch or debug tcr from
the IDE on the [src/testdata/java](../src/testdata/java) sample.

## Running and debugging another IDE

If you are using another IDE, use the following config:

- Working Directory: [src](../src)
- Command: `go run . -b ./testdata/java -w ./testdata/java -c ./testdata/java`

## Attaching a debugger to a running TCR instance

If you are using IntelliJ or GoLand, you can refer to
[this page](https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html)
for details on how to build and run a Go application and attach a debugger.

