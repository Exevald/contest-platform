import {createCefAdapter} from '../../common/cef/adapter'
import type {
	PlatformApi,
	StartupData,
} from '../types'

/**
 * Реализация CEF-адаптера для API генерации аудио.
 */
const cefCalls = {
	getStartupData: 'Platform_getStartupData',
}

const adapter = createCefAdapter()

const api: PlatformApi = {
	async getStartupData(): Promise<StartupData> {
		return adapter.call(cefCalls.getStartupData, [])
	},
}

export {
	api,
}
