Feature: Refactoring with TCR

  Scenario: Safe refactoring
    Given TCR is running
    And tests fully cover the behavior
    When I refactor the code
    Then tests pass
    And refactoring is committed automatically

  Scenario: Incomplete test coverage
    Given tests do not fully cover the code
    When I refactor
    Then TCR may revert my changes
    And I should improve test coverage first
