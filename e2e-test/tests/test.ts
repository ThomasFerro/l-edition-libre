import { Page, test as base, expect } from "@playwright/test";
import { Authentication } from "./authentication";
import { Manuscripts } from "./manuscripts";
import path from "path"

// TODO: Voir comment séparer les steps par domain proprement
class Given {
    constructor(private readonly authentication: Authentication) { }

    async IAmAnAuthenticatedWriter() {
        await this.authentication.authenticateAsWriter();
    }
};

class When {
    constructor(private readonly page: Page, private readonly manuscripts: Manuscripts) { }

    async ISubmitAManuscriptFor(manuscriptName: string) {
        await this.page.locator('[data-test-go-to="manuscripts"]').click()
        await this.page.locator('[data-test-new-manuscript-field="title"]').fill(this.manuscripts.get(manuscriptName));
        await this.page.locator('[data-test-new-manuscript-field="author"]').fill("Default author");

        const fileChooserPromise = this.page.waitForEvent("filechooser")
        await this.page.locator('[data-test-new-manuscript-field="file"]').click()
        const fileChooser = await fileChooserPromise
        await fileChooser.setFiles(path.join(__dirname, "assets/test.pdf"))

        await this.page.locator('[data-test="Submit new manuscript"]').click()
        await this.page.waitForLoadState("networkidle")
    }
}

class Then {
    constructor(private readonly page: Page, private readonly authentication: Authentication, private readonly manuscripts: Manuscripts) { }
    async TheFollowingManuscriptIsPendingReviewFromTheEditor(manuscriptName: string) {
        await this.page.goto("/manuscripts");
        const manuscript = this.page.locator('.manuscript', { hasText: this.manuscripts.get(manuscriptName) })
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
}

export const test = base.extend<{
    Authentication: Authentication,
    Manuscripts: Manuscripts,
    Given: Given,
    When: When,
    Then: Then
}>({
    Manuscripts: async ({}, use) => {
        const manuscripts = new Manuscripts();
        await use(manuscripts)
    },
    Authentication: async ({ page }, use) => {
        const authentication = new Authentication(page);
        await use(authentication);
    },
    Given: async ({ Authentication }, use) => {
        const given = new Given(Authentication);
        await use(given);
    },
    When: async ({ page, Manuscripts }, use) => {
        const when = new When(page, Manuscripts);
        await use(when);
    },
    Then: async ({ page, Authentication, Manuscripts }, use) => {
        const then = new Then(page, Authentication, Manuscripts);
        await use(then);
    },
}) 