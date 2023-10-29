import { test } from "./test";
test.describe('Submit a manuscript', () => {
  test('Submit a manuscript', async ({ Given, When, Then }) => {
    Given.IAmAnAuthenticatedWriter()
    When.ISubmitAManuscriptFor("My first novel")
    await Then.TheFollowingManuscriptIsPendingReviewFromTheEditor("My first novel")
  });
})

