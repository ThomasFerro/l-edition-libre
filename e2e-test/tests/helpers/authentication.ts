import { Page } from "@playwright/test";
import { Writers } from "./writers";
import { Navigation } from "./navigation";

export class Authentication {
    constructor(private readonly page: Page, private readonly navigation: Navigation) { }

    async givenIAmAnAuthenticatedWriter() {
        await this.authenticateAsWriter();
    }

    async givenIAmAuthenticatedAsAnotherWriter() {
        await this.authenticateAsWriter(Writers.AnotherAuthor)
    }

    async whenIAuthentifyAsAnEditor() {
        await this.authenticateAsEditor()
    }

    async authenticateAsWriter(writerName: string = Writers.FirstAuthor) {
        if (writerName == Writers.FirstAuthor) {
            return this.authenticate(process.env["AUTH0_WRITER_USERNAME"], process.env["AUTH0_WRITER_PASSWORD"])
        }
        if (writerName == Writers.AnotherAuthor) {
            return this.authenticate(process.env["AUTH0_SECOND_WRITER_USERNAME"], process.env["AUTH0_SECOND_WRITER_PASSWORD"])
        }
        throw new Error("Unknown writer " + writerName)
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
        await this.page.keyboard.press("Enter");
        await this.page.waitForSelector('[data-test="Disconnect"]')
    }
}
