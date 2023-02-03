@Writer @Manuscript
Feature: Submit a manusript
  As a writer,
  I want to submit manusripts
  In order to have them reviewed and eventually edited

  Background:
    Given I am an authentified writer


  @PDF
  Scenario: Submit a PDF manuscript
    When I submit a PDF manuscript for "My first novel"
    Then "My first novel" is pending review from the editor

  @PDF
  Scenario: Cancel a PDF manuscript submission
    Given I submitted a PDF manuscript for "My first novel"
    When I cancel the submission of "My first novel"
    Then "My first novel"'s submission is canceled

# TODO: Impossible d'annuler un manuscrit qui n'est pas pending