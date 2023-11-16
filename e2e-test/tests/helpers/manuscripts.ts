import { Page, expect } from "@playwright/test";
import { randomUUID } from "node:crypto";
import path from "path"
import { Authentication } from "./authentication";
import { WriterName } from "./writers";

export type ManuscriptName = string
export type ManuscriptUniqueIdentifier = string
export type ManuscriptStatus = "PendingReview" | "Canceled"

const cancelManuscriptSubmissionLocator = '[data-test-manuscript-action="Cancel"]';
export class Manuscripts {
    constructor(private readonly page: Page, private readonly authentication: Authentication) { }

    async givenISubmittedAManuscriptFor(manuscript: ManuscriptName) {
        await this.whenISubmitAManuscriptFor(manuscript)
    }

    async whenIGoToTheManuscriptsPage() {
        await this.goToManuscriptsPage()
    }

    async whenISubmitAManuscriptFor(manuscriptName: ManuscriptName) {
        await this.goToManuscriptsPage()
        await this.page.locator('[data-test-new-manuscript-field="title"]').fill(this.get(manuscriptName));
        await this.page.locator('[data-test-new-manuscript-field="author"]').fill("Default author");

        const fileChooserPromise = this.page.waitForEvent("filechooser")
        await this.page.locator('[data-test-new-manuscript-field="file"]').click()
        const fileChooser = await fileChooserPromise
        await fileChooser.setFiles(path.join(__dirname, "../assets/test.pdf"))

        await this.page.locator('[data-test="Submit new manuscript"]').click()
    }

    async thenTheFollowingManuscriptIsPendingReviewFromTheEditor(manuscriptName: ManuscriptName) {
        await this.manuscriptIsVisible(manuscriptName, "PendingReview");
        // TODO: Cette step = l'action de faire une review
        /*
        await this.authentication.authenticateAsEditor();
        await this.page.goto("/manuscripts/to-review");
        const manuscript = this.page.locator('.manuscript', { hasText: this.manuscripts.get(manuscriptName) })
        await expect(manuscript).toBeVisible()
        */
    }

    async thenMyManuscriptsAre(manuscripts: ManuscriptName[]) {
        for (const manuscript of manuscripts) {
            await this.thenTheFollowingManuscriptIsPendingReviewFromTheEditor(manuscript)
        }
    }

    async whenICancelSubmissionOfManuscript(manuscriptName: ManuscriptName) {
        const manuscript = this.manuscriptLocator(manuscriptName);
        await manuscript.locator(cancelManuscriptSubmissionLocator).click();
    }

    async thenSubmissionOfManuscriptIsCanceled(manuscriptName: ManuscriptName) {
        await this.manuscriptIsVisible(manuscriptName, "Canceled")
    }

    async givenSubmissionOfManuscriptWasCanceled(manuscriptName: ManuscriptName) {
        return this.whenICancelSubmissionOfManuscript(manuscriptName);
    }

    async thenICannotCancelSubmissionOfManuscript(manuscriptName: ManuscriptName) {
        const manuscript = this.manuscriptLocator(manuscriptName);
        await expect(manuscript.locator(cancelManuscriptSubmissionLocator)).not.toBeVisible();
    }

    async thenICannotSeeManuscript(manuscriptName: ManuscriptName) {
        const manuscript = this.manuscriptLocator(manuscriptName);
        await expect(manuscript.locator(cancelManuscriptSubmissionLocator)).not.toBeVisible();
    }

    async givenTheWriterSubmittedAManuscriptFor(author: WriterName, manuscriptName: ManuscriptName) {
        await this.authentication.authenticateAsWriter(author);
        await this.whenISubmitAManuscriptFor(manuscriptName);
    }

    async thenTheManuscriptsToReviewAre(manuscriptsToReview: { name: ManuscriptName; author: WriterName; }[]) {
        for (const manuscriptToReview of manuscriptsToReview) {
            await expect(this.page.locator(`[data-test-manuscript-to-review][data-test-manuscript-name="${manuscriptToReview.name}"][data-test-manuscript-author="${manuscriptToReview.author}"]`)).toBeVisible();
        }
        // TODO: Vérifier l'ordre
    }

    private manuscripts: Record<ManuscriptName, ManuscriptUniqueIdentifier> = {}
    private get(manuscriptName: ManuscriptName): ManuscriptUniqueIdentifier {
        let manuscriptIdentifier = this.manuscripts[manuscriptName]
        if (!manuscriptIdentifier) {
            manuscriptIdentifier = this.manuscripts[manuscriptName] = manuscriptName + randomUUID()
        }

        return manuscriptIdentifier
    }

    private manuscriptLocator(manuscriptName: ManuscriptName) {
        return this.page.locator('.manuscript', { hasText: this.get(manuscriptName) });
    }

    private async manuscriptIsVisible(manuscriptName: ManuscriptName, status: ManuscriptStatus) {
        const manuscript = this.manuscriptLocator(manuscriptName)
        await expect(manuscript).toBeVisible();
        await expect(manuscript.locator(`[data-test-manuscript-status="${status}"]`)).toBeVisible();
    }

    private async goToManuscriptsPage() {
        if (this.page.url().endsWith("manuscripts")) {
            return
        }
        await this.page.locator('[data-test-go-to="manuscripts"]').click()
    }
}
