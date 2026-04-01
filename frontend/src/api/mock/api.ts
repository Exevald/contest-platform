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
			languages: ['C++', 'Pascal', 'JavaScript', 'Python'],
			tabs: [
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
}

export {
	api,
}
