Feature: Limitations and trade-offs of TCR

  Scenario: Large changes
    Given I need to make a large structural change
    Then TCR may slow me down
    And I may need to temporarily disable it

  Scenario: Legacy code
    Given a legacy codebase with few tests
    Then TCR is hard to adopt
    And adding tests should be the first step

  Scenario: User frustration
    Given frequent reverts
    Then TCR can feel frustrating
    But this frustration highlights design or testing issues
