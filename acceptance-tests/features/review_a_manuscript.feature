@Editor @Manuscript
Feature: Review a manusript
  As a edit,
  I want to review manusripts
  In order to have them published

  Scenario: List manuscripts to be reviewed, by order of submission
    Given the writer "First author" submitted a manuscript for "My first novel"
    And the writer "Another author" submitted a manuscript for "Another novel"
    When I authentify as an editor
    Then the manuscripts to be reviewed are the following
      | Title          | Author         |
      | My first novel | First author   |
      | Another novel  | Another author |

  Scenario: Review a manuscript
    Given a writer submitted a manuscript for "My first novel"
    And I am an authentified editor
    When I review positively the manuscript for "My first novel"
    Then "My first novel" is eventually published

  @Error
  Scenario: Only manuscripts pending review can be reviewed
    Given a writer submitted a manuscript for "My first novel"
    And submission of "My first novel" was canceled
    And I am an authentified editor
    When I review positively the manuscript for "My first novel"
    Then the error "AManuscriptShouldBePendingReviewToBeReviewed" is thrown

  @Error
  Scenario: Manuscripts can only be reviewed by an editor
    Given I am an authentified writer
    And I submitted a manuscript for "My first novel"
    When I review positively the manuscript for "My first novel"
    Then the error "ManuscriptNotFound" is thrown

# TODO: ManuscriptNeedsRework ?
# TODO: Ouvrir le document