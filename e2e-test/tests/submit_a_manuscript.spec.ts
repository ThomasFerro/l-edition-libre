import { test } from "./helpers/test";
test.describe('Submit a manuscript', () => {
  test('Submit a manuscript', async ({ Manuscripts, Authentication }) => {
    await Authentication.givenIAmAnAuthenticatedWriter()
    await Manuscripts.whenISubmitAManuscriptFor("My first novel")
    await Manuscripts.thenTheFollowingManuscriptIsPendingReviewFromTheEditor("My first novel")
  });
})

