Feature: Onboarding a new TCR user

  Scenario: Prerequisites for using TCR
    Given I have a project under version control
    And the project has a reliable automated test suite
    Then I am eligible to use TCR

  Scenario: First encounter with TCR
    Given I have never used TCR before
    When I start TCR for the first time
    Then I should quickly understand that tests control commits
    And failures cause an automatic revert

  Scenario: Learning expectations
    Given I am a new user
    Then I should expect a strict and unforgiving workflow
    And I should expect fast feedback
