import {sleep} from '@reatom/core'
import type {
	GetDataArgs,
	PlatformApi, ResetTaskArgs, SendFileArgs, StartupData,
} from '../types'
import {newTaskData, taskData} from './data'

const api: PlatformApi = {
	async getStartupData(): Promise<StartupData> {
		return {
			title: 'Rocket Pizza',
			languages: [
				{
					name: 'C++',
					extension: 'cpp'
				},
			],
			tasks: [
				{
					id: '1',
					type: 'table',
					label: 'Фильтр заказов',
				},
				{
					id: '2',
					type: 'table',
					label: 'Карта спроса',
				},

				{
					id: '3',
					type: 'table',
					label: 'Генератор маршрутов',
				},
			],
			workspaceViews: [
				{
					id: 'statement',
					label: 'Данные',
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
}

export {
	api,
}
