import { Page, test as base, expect } from "@playwright/test";
import { Authentication } from "./authentication";

// TODO: Voir comment séparer les steps par domain proprement
class Given {
    constructor(private readonly authentication: Authentication) { }

    async IAmAnAuthenticatedWriter() {
        await this.authentication.authenticateAsWriter();
    }
};
class When {
    constructor(private readonly options: { baseURL }) { }

    ISubmitAManuscriptFor(manuscriptName: string) { console.log(this.options) }
}
class Then {
    constructor(private readonly page: Page, private readonly authentication: Authentication) { }
    async TheFollowingManuscriptIsPendingReviewFromTheEditor(manuscriptName: string) {
        await this.authentication.authenticateAsEditor();
        await this.page.goto("/manuscripts/to-review")
        const manuscript = this.page.locator('.manuscript', { hasText: manuscriptName })
        await expect(manuscript).toBeVisible()
    }
}

export const test = base.extend<{
    Authentication: Authentication,
    Given: Given,
    When: When,
    Then: Then
}>({
    Authentication: async ({page}, use ) => {
        const authentication = new Authentication(page);
        await use(authentication);
    },
    Given: async ({ Authentication }, use) => {
        const given = new Given(Authentication);
        await use(given);
    },
    When: async ({ page, baseURL }, use) => {
        const when = new When({ baseURL });
        await use(when);
    },
    Then: async ({ page, Authentication }, use) => {
        const then = new Then(page, Authentication);
        await use(then);
    },
}) 