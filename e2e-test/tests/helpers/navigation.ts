import { Page } from "@playwright/test";
import { ApplicationPage } from "./test";

export class Navigation {
    constructor(private readonly page: Page) { }

    navigateTo(applicationPage: ApplicationPage) {
        // TODO: Ok avec des query params ? 
        if (this.page.url().endsWith(applicationPage)) {
            return
        }
        await this.page.waitForURL(url => url.pathname.endsWith(applicationPage));
    }
}
