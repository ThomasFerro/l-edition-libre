import { Page, test as base, expect } from "@playwright/test";

class Given {
    IAmAnAuthentifiedWriter() { }
};
class When {
    constructor(private readonly options: { baseURL }) { }

    ISubmitAManuscriptFor(manuscriptName: string) { console.log(this.options) }
}
class Then {
    constructor(private readonly page: Page) { }
    async TheFollowingManuscriptIsPendingReviewFromTheEditor(manuscriptName: string) {
        // TODO: Connecter en tant qu'éditeur
        await this.page.goto("")
        const manuscript = this.page.locator('.manuscript', { hasText: manuscriptName })
        await expect(manuscript).toBeVisible()
    }
}

export const test = base.extend<{
    Given: Given,
    When: When,
    Then: Then
}>({
    Given: async ({ page }, use) => {
        const given = new Given();
        await use(given);
    },
    When: async ({ page, baseURL }, use) => {
        const when = new When({ baseURL });
        await use(when);
    },
    Then: async ({ page }, use) => {
        const then = new Then(page);
        await use(then);
    },
}) 