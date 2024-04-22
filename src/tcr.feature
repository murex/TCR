Feature: Mobbing with TCR work session

  Scenario: Mobbing with TCR and test passing
    Given a TCR engine is started
    When a developer writes some code
    And the developer runs the tests
    Then if the tests pass, the code is committed

  Scenario: Mobbing with TCR and test failing
    Given a TCR engine is started
    When a developer writes some code
    And the developer runs the tests
    Then if the tests fail, the code is reverted

  Scenario: Mobbing with TCR, and switching role when timer expires
    Given a TCR engine is started
    When the mob timer expires
    Then the developer switches roles with another developer in the mob