Feature: Basic TCR workflow

  Scenario: Writing code that passes tests
    Given TCR is running
    When I write code that satisfies the tests
    Then tests pass
    And my changes are committed automatically

  Scenario: Writing code that breaks tests
    Given TCR is running
    When I write code that causes a test to fail
    Then tests fail
    And my changes are reverted automatically

  Scenario: Small incremental changes
    Given TCR is running
    Then I am encouraged to make very small changes
