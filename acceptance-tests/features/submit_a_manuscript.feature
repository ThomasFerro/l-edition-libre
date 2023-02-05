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
    Then submission of "My first novel" is canceled

  @Error
  Scenario: A manuscript should be pending review for its submission to be canceled
    Given I submitted a PDF manuscript for "My first novel"
    And submission of "My first novel" was canceled
    When I cancel the submission of "My first novel"
    Then the error "AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled" is thrown

# TODO: Impossible d'annuler un manuscrit qui n'est pas pending
# TODO: Impossible de voir le status d'un manuscript ou en annuler la soumission si ce n'est pas le notre
# TODO: Véritablement téléverser un document (vérifier qu'on l'a bien persisté en le récupérant ?)