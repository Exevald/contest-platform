type Task = {
	id: string,
	type: TabType,
	label: string,
}

type Language = {
	name: string,
	extension: string,
}

type TabType = 'table'
type WorkspaceScreen = 'statement' | 'submission_history'
type WorkspaceView = {
	id: WorkspaceScreen,
	label: string,
}

export type {
	Task,
	Language,
	WorkspaceScreen,
	WorkspaceView,
}
