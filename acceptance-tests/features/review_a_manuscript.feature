@Editor @Manuscript
Feature: Review a manusript
  As a edit,
  I want to review manusripts
  In order to have them published

  Scenario: Review a manuscript
    Given a writer submitted a manuscript for "My first novel"
    And I am an authentified editor
    When I review positively the manuscript for "My first novel"
    Then "My first novel" is eventually published

  @Error
  Scenario: Only the manuscripts pending review can be reviewed
    Given a writer submitted a manuscript for "My first novel"
    And submission of "My first novel" was canceled
    And I am an authentified editor
    When I review positively the manuscript for "My first novel"
    Then the error "AManuscriptShouldBePendingReviewToBeReviewed" is thrown

# TODO: ManuscriptNeedsRework ?
# TODO: Gestion des droits, un écrivain ne peut pas review
# TODO: Liste des pending review + ouvrir le document