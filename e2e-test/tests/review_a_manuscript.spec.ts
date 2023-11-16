import { test } from "./helpers/test";
import { Writers } from "./helpers/writers";

test.describe('Review a manuscript', () => {
  test('List manuscripts to be reviewed', async ({ Manuscripts, Authentication }) => {
    await Manuscripts.givenTheWriterSubmittedAManuscriptFor(Writers.FirstAuthor, "My first novel")
    await Manuscripts.givenTheWriterSubmittedAManuscriptFor(Writers.AnotherAuthor, "Another novel")
    await Authentication.whenIAuthentifyAsAnEditor()
    await Manuscripts.thenTheManuscriptsToReviewAre([{
      name: 'My first novel',
      author: Writers.FirstAuthor
    }, {
      name: 'Another novel',
      author: Writers.AnotherAuthor
    }])
  });
})
