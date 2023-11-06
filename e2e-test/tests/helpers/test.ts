import { Page, test as base, expect } from "@playwright/test";
import { Authentication } from "./authentication";
import { Manuscripts } from "./manuscripts";

export const test = base.extend<{
    Authentication: Authentication,
    Manuscripts: Manuscripts,
}>({
    Manuscripts: async ({ page, Authentication }, use) => {
        const manuscripts = new Manuscripts(page, Authentication);
        await use(manuscripts)
    },
    Authentication: async ({ page }, use) => {
        const authentication = new Authentication(page);
        await use(authentication);
    },
}) 