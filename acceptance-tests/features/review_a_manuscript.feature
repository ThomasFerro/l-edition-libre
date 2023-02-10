@Editor @Manuscript
Feature: Review a manusript
  As a edit,
  I want to review manusripts
  In order to have them published

  Background:
    Given a writer submitted a manuscript for "My first novel"
    And I am an authentified editor

  # Scenario: Review a manuscript
  #   When I review positively the manuscript for "My first novel"
  #   Then "My first novel" is eventually published

# TODO: ManuscriptNeedsRework ?
# TODO: Gestion des droits, un écrivain ne peut pas review
# TODO: Liste des pending review + ouvrir le document
# TODO: Impossible de review des manuscrits qui ne sont pas readyToReview