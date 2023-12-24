import { Page } from "@playwright/test";

export type ApplicationPage = '/manuscripts' | '/manuscripts-to-review';
export class Navigation {
    constructor(private readonly page: Page) { }

    async navigateTo(applicationPage: ApplicationPage) {
        // TODO: Ok avec des query params ? 
        if (this.page.url().endsWith(applicationPage)) {
            return
        }

        await this.page.locator(`[data-test-go-to="${applicationPage}"]`).click()
        await this.page.waitForURL(url => url.pathname.endsWith(applicationPage));
    }
}
