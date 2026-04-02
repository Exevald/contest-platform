import type {
	StartupApi,
} from '../common/startApp/api/StartupApi'
import type {Language, Task, ThemeKey, WorkspaceView} from '../model/types'

type StartupData = {
	title: string,
	languages: Language[],
	tasks: Task[],
	workspaceViews: WorkspaceView[],
	participantCode: string,
	selectedTheme: string,
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

type StartSessionArgs = {
	participantCode: string,
	theme: ThemeKey,
}

type SubmissionStatus = {
	submissionId: string,
	problemId: string,
	language: string,
	verdict: string,
	createdAt: string,
	compilationOutput: string,
}

type PlatformApi = StartupApi<StartupData> & {
	getData: <T>(data: GetDataArgs) => Promise<T>,
	sendFile: <T>(data: SendFileArgs) => Promise<T>,
	resetTask: <T>(data: ResetTaskArgs) => Promise<T>,
	getLatestSubmission: <T>(data: SubmissionArgs) => Promise<T>,
	getSubmissionStatus: <T>(data: SubmissionArgs) => Promise<T>,
	getSubmissionHistory: <T>(data: SubmissionArgs) => Promise<T>,
	startSession: (data: StartSessionArgs) => Promise<StartupData>,
}

export {
	PlatformApi,
	StartupData,
	GetDataArgs,
	SendFileArgs,
	ResetTaskArgs,
	SubmissionArgs,
	SubmissionStatus,
	StartSessionArgs,
}
