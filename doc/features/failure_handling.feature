Feature: Failure handling in TCR

  Scenario: Test failure
    Given TCR is running
    When a test fails
    Then all uncommitted changes are reverted

  Scenario: Compilation failure
    Given TCR is running
    When the code does not compile
    Then changes are reverted

  Scenario: Non-deterministic tests
    Given tests are flaky
    Then using TCR becomes dangerous
    And tests should be stabilized before continuing
