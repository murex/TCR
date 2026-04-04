Feature: TCR overview

  TCR stands for "Test && Commit || Revert".

  It enforces a strict development loop where:
  - Tests are always run
  - Passing tests result in a commit
  - Failing tests result in a revert

  Scenario: Understanding the TCR feedback loop
    Given I am working on a codebase with automated tests
    When I make a change to the code
    Then tests are automatically executed
    And if tests pass, my changes are committed
    And if tests fail, my changes are reverted

  Scenario: No intermediate broken state
    Given I am using TCR
    When I introduce a failing change
    Then the codebase never remains in a broken state
