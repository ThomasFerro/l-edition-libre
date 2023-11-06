import { test } from "./helpers/test";
import { Writers } from "./helpers/writers";
test.describe('Submit a manuscript', () => {
  test('Submit a manuscript', async ({ Manuscripts, Authentication }) => {
    await Authentication.givenIAmAnAuthenticatedWriter()
    await Manuscripts.whenISubmitAManuscriptFor("My first novel")
    await Manuscripts.thenTheFollowingManuscriptIsPendingReviewFromTheEditor("My first novel")
  });

  test('List submitted manuscripts', async ({ Manuscripts, Authentication }) => {
    await Authentication.givenIAmAnAuthenticatedWriter()
    await Manuscripts.givenISubmittedAManuscriptFor("My first novel")
    await Manuscripts.givenISubmittedAManuscriptFor("My second novel")
    await Authentication.givenIAmAuthenticatedAsAnotherWriter()
    await Manuscripts.whenISubmitAManuscriptFor("Essay #1")
    await Manuscripts.whenISubmitAManuscriptFor("Essay #2")
    await Manuscripts.thenMyManuscriptsAre(["Essay #1", "Essay #2"])
  });
})