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
				{
					name: 'Pascal',
					extension: 'pas'
				},
				{
					name: 'JavaScript',
					extension: 'js'
				},
				{
					name: 'Python',
					extension: 'py'
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
