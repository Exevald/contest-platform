import type {
	StartupApi,
} from '../common/startApp/api/StartupApi'
import type {Task} from '../model/types'

type StartupData = {
	title: string,
	languages: string[],
	tabs: Task[],
}

type GetDataArgs = {
	id: string
}

type SendFileArgs = {
	id: string
	language: string
	text: string
}

type ResetTaskArgs = {
	id: string
}

type PlatformApi = StartupApi<StartupData> & {
	getData: <T>(data: GetDataArgs) => Promise<T>,
	sendFile: <T>(data: SendFileArgs) => Promise<T>,
	resetTask: <T>(data: ResetTaskArgs) => Promise<T>,
}

export {
	PlatformApi,
	StartupData,
	GetDataArgs,
	SendFileArgs,
	ResetTaskArgs,
}
