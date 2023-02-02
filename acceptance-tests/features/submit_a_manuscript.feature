@Writer @Manuscript
Feature: Submit a manusript
  As a writer,
  I want to submit manusripts
  In order to have them reviewed and eventually edited

  @PDF
  Scenario: Submit a PDF manuscript
    Given I am an authentified writer
    When I submit a PDF manuscript for "My first novel"
    Then "My first novel" is pending review from the editor