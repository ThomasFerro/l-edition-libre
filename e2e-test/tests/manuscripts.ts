import { randomUUID } from "node:crypto";

export type ManuscriptName = string
export type ManuscriptUniqueIdentifier = string

export class Manuscripts {
    private manuscripts: Record<ManuscriptName, ManuscriptUniqueIdentifier> = {}

    get(manuscriptName: ManuscriptName): ManuscriptUniqueIdentifier {
        let manuscriptIdentifier = this.manuscripts[manuscriptName]
        if (!manuscriptIdentifier) {
            manuscriptIdentifier = this.manuscripts[manuscriptName] = manuscriptName + randomUUID()
        }

        return manuscriptIdentifier 
    }
}