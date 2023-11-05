import { Page, expect } from "@playwright/test";
import { randomUUID } from "node:crypto";
import path from "path"

export type ManuscriptName = string
export type ManuscriptUniqueIdentifier = string

export class Manuscripts {
    private manuscripts: Record<ManuscriptName, ManuscriptUniqueIdentifier> = {}

    constructor(private readonly page: Page) { }

    async whenISubmitAManuscriptFor(manuscriptName: string) {
        await this.page.locator('[data-test-go-to="manuscripts"]').click()
        await this.page.locator('[data-test-new-manuscript-field="title"]').fill(this.get(manuscriptName));
        await this.page.locator('[data-test-new-manuscript-field="author"]').fill("Default author");

        const fileChooserPromise = this.page.waitForEvent("filechooser")
        await this.page.locator('[data-test-new-manuscript-field="file"]').click()
        const fileChooser = await fileChooserPromise
        await fileChooser.setFiles(path.join(__dirname, "../assets/test.pdf"))

        await this.page.locator('[data-test="Submit new manuscript"]').click()
        await this.page.waitForLoadState("networkidle")
    }

    async thenTheFollowingManuscriptIsPendingReviewFromTheEditor(manuscriptName: string) {
        await this.page.goto("/manuscripts");
        const manuscript = this.page.locator('.manuscript', { hasText: this.get(manuscriptName) })
        await expect(manuscript).toBeVisible()
        await expect(manuscript.locator('[data-test-manuscript-status="PendingReview"]')).toBeVisible()
        // TODO: Cette step = l'action de faire une review
        /*
        await this.authentication.authenticateAsEditor();
        await this.page.goto("/manuscripts/to-review");
        const manuscript = this.page.locator('.manuscript', { hasText: this.manuscripts.get(manuscriptName) })
        await expect(manuscript).toBeVisible()
        */
    }

    get(manuscriptName: ManuscriptName): ManuscriptUniqueIdentifier {
        let manuscriptIdentifier = this.manuscripts[manuscriptName]
        if (!manuscriptIdentifier) {
            manuscriptIdentifier = this.manuscripts[manuscriptName] = manuscriptName + randomUUID()
        }

        return manuscriptIdentifier
    }
}