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

  @Error @Users
  Scenario: Only the writer of a manuscript see its submission status
    Given I submitted a PDF manuscript for "My first novel"
    And I am authentified as another writer
    When I try to get the submission status of "My first novel"
    Then the error "ManuscriptNotFound" is thrown

  @Error @Users
  Scenario: Only the writer of a manuscript can cancel its submission
    Given I submitted a PDF manuscript for "My first novel"
    And I am authentified as another writer
    When I cancel the submission of "My first novel"
    Then the error "ManuscriptNotFound" is thrown

# TODO: Impossible d'annuler un manuscrit qui n'est pas pending
# TODO: Véritablement téléverser un document (vérifier qu'on l'a bien persisté en le récupérant ?)