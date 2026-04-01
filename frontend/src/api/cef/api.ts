import type {GetDataArgs, PlatformApi, ResetTaskArgs, SendFileArgs, StartupData, SubmissionArgs,} from '../types'

type WailsAppApi = {
    GetStartupData: () => Promise<StartupData>,
    GetData: (id: string) => Promise<string>,
    SendFile: (id: string, language: string, text: string) => Promise<unknown>,
    ResetTask: (id: string) => Promise<string>,
    GetLatestSubmission: (id: string) => Promise<unknown>,
    GetSubmissionStatus: (id: string) => Promise<unknown>,
}

const getBackend = (): WailsAppApi => {
    const backend = (window as any)?.go?.main?.App

    if (!backend) {
        throw new Error('Wails backend is not available')
    }

    return backend as WailsAppApi
}

const api: PlatformApi = {
    async getStartupData(): Promise<StartupData> {
        return getBackend().GetStartupData()
    },

    async getData({id}: GetDataArgs): Promise<string> {
        return getBackend().GetData(id)
    },

    async sendFile({id, language, text}: SendFileArgs): Promise<unknown> {
        return getBackend().SendFile(id, language, text)
    },

    async resetTask({id}: ResetTaskArgs): Promise<string> {
        return getBackend().ResetTask(id)
    },

    async getLatestSubmission({id}: SubmissionArgs): Promise<unknown> {
        return getBackend().GetLatestSubmission(id)
    },

    async getSubmissionStatus({id}: SubmissionArgs): Promise<unknown> {
        return getBackend().GetSubmissionStatus(id)
    },
}

export {
    api,
}
