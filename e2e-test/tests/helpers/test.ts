import { test as base } from "@playwright/test";
import { Authentication } from "./authentication";
import { Manuscripts } from "./manuscripts";

export const test = base.extend<{
    Authentication: Authentication,
    Manuscripts: Manuscripts,
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
    Manuscripts: async ({ page, Authentication }, use) => {
        const manuscripts = new Manuscripts(page, Authentication);
        await use(manuscripts)
    },
    Authentication: async ({ page }, use) => {
        const authentication = new Authentication(page);
        await use(authentication);
    },
}) 