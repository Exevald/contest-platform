import {sleep} from '@reatom/core'
import type {GetDataArgs, PlatformApi, ResetTaskArgs, SendFileArgs, StartupData,} from '../types'
import {getAllStartupTasks} from '../../config/theme'
import {newTaskData, taskData} from './data'

let participantCode = ''
let selectedTheme = ''

const api: PlatformApi = {
    async getStartupData(): Promise<StartupData> {
        return {
            title: 'ContestPlatform',
            languages: [
                {
                    name: 'C++',
                    extension: 'cpp'
                },
                {
                    name: 'Go',
                    extension: 'go'
                },
                {
                    name: 'JavaScript',
                    extension: 'js'
                },
                {
                    name: 'Pascal',
                    extension: 'pas'
                },
                {
                    name: 'PHP',
                    extension: 'php'
                },
                {
                    name: 'Python',
                    extension: 'py'
                },
            ],
            tasks: getAllStartupTasks(),
            participantCode,
            selectedTheme,
            workspaceViews: [
                {
                    id: 'statement',
                    label: 'Задание',
                },
                {
                    id: 'submission_history',
                    label: 'История посылок',
                },
            ],
        }
    },

    async getData(data: GetDataArgs) {
        await sleep(2000)
        return taskData[data.id]
    },

    sendFile(data: SendFileArgs) {
        console.log(data)
        return newTaskData[data.id]
    },

    resetTask(data: ResetTaskArgs) {
        return taskData[data.id]
    },

    getLatestSubmission() {
        return Promise.resolve(null)
    },

    getSubmissionStatus() {
        return Promise.resolve(null)
    },

    getSubmissionHistory() {
        return Promise.resolve([])
    },

    startSession(data) {
        participantCode = data.participantCode
        selectedTheme = data.theme
        return this.getStartupData()
    },
}

export {
    api,
}
