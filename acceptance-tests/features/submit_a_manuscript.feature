@Writer @Manuscript
Feature: Submit a manusript
  As a writer,
  I want to submit manusripts
  In order to have them reviewed and eventually edited

  Background:
    Given I am an authenticated writer

  Scenario: Submit a manuscript
    When I submit a manuscript for "My first novel"
    Then "My first novel" is pending review from the editor

  Scenario: List submitted manuscripts
    Given the writer "First author" submitted a manuscript for "My first novel"
    And the writer "First author" submitted a manuscript for "My second novel"
    And I am authenticated as another writer
    When I submit a manuscript for "Essay #1"
    And I submit a manuscript for "Essay #2"
    Then my manuscripts are the following
      | Title    |
      | Essay #1 |
      | Essay #2 |

  Scenario: Cancel a manuscript submission
    Given I submitted a manuscript for "My first novel"
    When I cancel the submission of "My first novel"
    Then submission of "My first novel" is canceled

  @Error
  Scenario: A manuscript should be pending review for its submission to be canceled
    Given I submitted a manuscript for "My first novel"
    And submission of "My first novel" was canceled
    When I cancel the submission of "My first novel"
    Then the error "AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled" is thrown

  @Error @Users
  Scenario: Only the writer of a manuscript can see its submission status
    Given I submitted a manuscript for "My first novel"
    And I am authenticated as another writer
    When I try to get the submission status of "My first novel"
    Then the error "ManuscriptNotFound" is thrown

  @Error @Users
  Scenario: Only the writer of a manuscript can cancel its submission
    Given I submitted a manuscript for "My first novel"
    And I am authenticated as another writer
    When I cancel the submission of "My first novel"
    Then the error "ManuscriptNotFound" is thrown