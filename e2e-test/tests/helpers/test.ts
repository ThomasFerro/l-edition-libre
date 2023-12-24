import { test as base } from "@playwright/test";
import { Authentication } from "./authentication";
import { Manuscripts } from "./manuscripts";
import { Navigation } from "./navigation";

export const test = base.extend<{
    Authentication: Authentication,
    Manuscripts: Manuscripts,
    Navigation: Navigation
}>({
    page: async ({ page }, use) => {
        page.on('console', (consoleEvent) => {
            console.log(`[${consoleEvent.type()}]${consoleEvent.text()}`)
        })
        page.on('requestfailed', (request) => {
            console.log(`[FAIL ${request.method()} ${request.url()}]${request.failure()?.errorText}`)
        })
        await use(page)
    },
    Navigation: async ({ page }, use) => {
        const navigation = new Navigation(page);
        await use(navigation)
    },
    Manuscripts: async ({ page, Authentication, Navigation }, use) => {
        const manuscripts = new Manuscripts(page, Authentication, Navigation);
        await use(manuscripts)
    },
    Authentication: async ({ page, Navigation }, use) => {
        const authentication = new Authentication(page, Navigation);
        await use(authentication);
    },
}) 
