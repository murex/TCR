Feature: Mob Timer

  Scenario: Check timer duration trace at startup
    Given a TCR engine with a mob timer enabled
    When the TCR engine starts
    Then the timer duration should be traced exactly once

  Scenario: Mob timer should be off in solo mode
    Given a TCR engine with a mob timer enabled
    And the TCR engine is set to solo mode
    When the TCR engine starts
    Then the timer duration should not be traced at all