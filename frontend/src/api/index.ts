import type {PlatformApi} from './types'
import {api as backendApi} from './cef/api'
import {api as mockApi} from './mock/api'

function hasBackend() {
    return Boolean((window as any)?.go?.main?.App)
}

const api: PlatformApi = {
    getStartupData() {
        return (hasBackend() ? backendApi : mockApi).getStartupData()
    },

    getData(args) {
        return (hasBackend() ? backendApi : mockApi).getData(args)
    },

    sendFile(args) {
        return (hasBackend() ? backendApi : mockApi).sendFile(args)
    },

    resetTask(args) {
        return (hasBackend() ? backendApi : mockApi).resetTask(args)
    },

    getLatestSubmission(args) {
        return (hasBackend() ? backendApi : mockApi).getLatestSubmission(args)
    },

    getSubmissionStatus(args) {
        return (hasBackend() ? backendApi : mockApi).getSubmissionStatus(args)
    },

    getSubmissionHistory(args) {
        return (hasBackend() ? backendApi : mockApi).getSubmissionHistory(args)
    },

    startSession(args) {
        return (hasBackend() ? backendApi : mockApi).startSession(args)
    },
}

export {
    api,
}
