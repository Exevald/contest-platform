import type {
	StartupApi,
} from '../common/startApp/api/StartupApi'
import type {Language, Task} from '../model/types'

type StartupData = {
	title: string,
	languages: Language[],
	tasks: Task[],
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

type SubmissionArgs = {
	id: string
}

type SubmissionStatus = {
	submissionId: string,
	problemId: string,
	language: string,
	verdict: string,
	createdAt: string,
	testsPassed: number,
	testsTotal: number,
}

type PlatformApi = StartupApi<StartupData> & {
	getData: <T>(data: GetDataArgs) => Promise<T>,
	sendFile: <T>(data: SendFileArgs) => Promise<T>,
	resetTask: <T>(data: ResetTaskArgs) => Promise<T>,
	getLatestSubmission: <T>(data: SubmissionArgs) => Promise<T>,
	getSubmissionStatus: <T>(data: SubmissionArgs) => Promise<T>,
}

export {
	PlatformApi,
	StartupData,
	GetDataArgs,
	SendFileArgs,
	ResetTaskArgs,
	SubmissionArgs,
	SubmissionStatus,
}
