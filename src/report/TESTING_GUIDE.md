# Reporter Testing Guide

This guide explains how to write tests that use the TCR reporter system in an isolated manner to prevent test interference.

## Problem

Previously, all tests using the reporter system shared a singleton reporter instance, which caused test interference. Messages from one test could leak into another test, leading to flaky or incorrect test results.

## Solution

We have refactored the reporter system to support isolated reporter instances for testing. Each test now uses its own reporter instance, ensuring complete isolation.

## How to Use Isolated Reporters in Tests

### Basic Pattern

Use `TestWithIsolatedReporter` for tests that need to capture and verify reporter messages:

```go
func Test_my_function_reports_messages(t *testing.T) {
    report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
        // Your test code here that generates messages
        myFunction() // This should call report.PostInfo(), etc.
        
        // Stop the sniffer and verify results
        sniffer.Stop()
        assert.Equal(t, 1, sniffer.GetMatchCount())
        assert.Equal(t, "expected message", sniffer.GetAllMatches()[0].Payload.ToString())
    })
}
```

### Pattern with Message Filters

Use `TestWithIsolatedReporterAndFilters` when you only want to capture specific types of messages:

```go
func Test_my_function_reports_errors(t *testing.T) {
    report.TestWithIsolatedReporterAndFilters(func(reporter *report.Reporter, sniffer *report.Sniffer) {
        // Your test code here
        myFunctionThatMayGenerateErrors()
        
        sniffer.Stop()
        assert.Equal(t, 1, sniffer.GetMatchCount()) // Only error messages counted
    }, func(msg report.Message) bool {
        return msg.Type.Category == report.Error
    })
}
```

### Multiple Filters

You can provide multiple filters - the sniffer will capture messages that match ANY of the filters:

```go
func Test_my_function_reports_warnings_and_errors(t *testing.T) {
    report.TestWithIsolatedReporterAndFilters(func(reporter *report.Reporter, sniffer *report.Sniffer) {
        myFunction()
        sniffer.Stop()
        // Will capture both warning and error messages
    }, 
    func(msg report.Message) bool { return msg.Type.Category == report.Warning },
    func(msg report.Message) bool { return msg.Type.Category == report.Error })
}
```

## Key Benefits

1. **Test Isolation**: Each test uses its own reporter instance
2. **No Interference**: Messages from one test cannot affect another test
3. **Automatic Cleanup**: The original default reporter is automatically restored
4. **Exception Safety**: Cleanup happens even if the test panics

## What Changed

### Before (Old Pattern - DON'T USE)
```go
func Test_old_pattern(t *testing.T) {
    sniffer := report.NewSniffer()
    myFunction() // Uses global reporter
    sniffer.Stop()
    assert.Equal(t, 1, sniffer.GetMatchCount())
}
```

### After (New Pattern - USE THIS)
```go
func Test_new_pattern(t *testing.T) {
    report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
        myFunction() // Uses isolated reporter
        sniffer.Stop()
        assert.Equal(t, 1, sniffer.GetMatchCount())
    })
}
```

## Available Message Types

The reporter supports these message categories:

- `report.Normal` - Regular messages
- `report.Info` - Information messages
- `report.Title` - Title messages
- `report.Success` - Success messages
- `report.Warning` - Warning messages
- `report.Error` - Error messages
- `report.RoleEvent` - Role-related events (driver/navigator)
- `report.TimerEvent` - Timer-related events

## Message Emphasis

Messages can have emphasis (boolean flag) for highlighting important information:

```go
// Check if a message has emphasis
if msg.Type.Emphasis {
    // This message should be highlighted
}
```

## Common Patterns

### Testing Message Content
```go
report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
    myFunction()
    sniffer.Stop()
    
    messages := sniffer.GetAllMatches()
    assert.Equal(t, "Expected message text", messages[0].Payload.ToString())
    assert.Equal(t, report.Info, messages[0].Type.Category)
    assert.True(t, messages[0].Type.Emphasis)
})
```

### Testing Message Count
```go
report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
    myFunction()
    sniffer.Stop()
    
    assert.Equal(t, 3, sniffer.GetMatchCount()) // Expecting exactly 3 messages
})
```

### Testing Specific Message Types
```go
report.TestWithIsolatedReporterAndFilters(func(reporter *report.Reporter, sniffer *report.Sniffer) {
    myFunction()
    sniffer.Stop()
    
    assert.Equal(t, 1, sniffer.GetMatchCount()) // Only warnings captured
}, func(msg report.Message) bool {
    return msg.Type.Category == report.Warning
})
```

## Migration Guide

If you have existing tests using the old pattern:

1. Replace `sniffer := report.NewSniffer()` with `report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {`
2. Move your test logic inside the function
3. Add the closing `})` 
4. If you were using filters, use `TestWithIsolatedReporterAndFilters` instead
5. Remove manual `sniffer.Stop()` calls if they're at the end (cleanup is automatic)

## Advanced Usage

### Custom Reporter Instance

If you need direct access to the reporter instance:

```go
report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
    // You can post messages directly to the isolated reporter
    reporter.PostInfo("test message")
    reporter.PostErrorWithEmphasis("error with emphasis")
    
    sniffer.Stop()
    assert.Equal(t, 2, sniffer.GetMatchCount())
})
```

### Testing Components That Use Global Report Functions

Most functions in the codebase use global report functions like `report.PostInfo()`. These automatically use the current default reporter, which is temporarily replaced during the test:

```go
// This function uses report.PostInfo() internally
func myBusinessLogic() {
    report.PostInfo("Processing started")
    // ... business logic ...
    report.PostInfo("Processing completed")
}

func Test_business_logic_reports_progress(t *testing.T) {
    report.TestWithIsolatedReporter(func(reporter *report.Reporter, sniffer *report.Sniffer) {
        myBusinessLogic() // Will use the isolated reporter automatically
        sniffer.Stop()
        
        assert.Equal(t, 2, sniffer.GetMatchCount())
        assert.Equal(t, "Processing started", sniffer.GetAllMatches()[0].Payload.ToString())
        assert.Equal(t, "Processing completed", sniffer.GetAllMatches()[1].Payload.ToString())
    })
}
```

This ensures that your tests are completely isolated while testing real application behavior.