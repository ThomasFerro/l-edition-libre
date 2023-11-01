import { Page } from "@playwright/test";

export class Authentication {
    constructor(private readonly page: Page) { }

    async authenticateAsWriter() {
        await this.authenticate(process.env["AUTH0_WRITER_USERNAME"], process.env["AUTH0_WRITER_PASSWORD"])
    }

    async authenticateAsEditor() {
        await this.authenticate(process.env["AUTH0_EDITOR_USERNAME"], process.env["AUTH0_EDITOR_PASSWORD"])
    }

    private async authenticate(login: string, password: string) {
        await this.page.goto("");
        const disconnectButton = this.page.locator('[data-test="Disconnect"]')
        if (await disconnectButton.isVisible()) {
            await disconnectButton.click()
        }
        await this.page.locator('[data-test="Go to connection page"]').click()
        await this.page.locator("#username").fill(login)
        await this.page.locator("#password").fill(password)
        await this.page.keyboard.press("Enter")
    }
}