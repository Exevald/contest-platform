type Task = {
	id: string,
	type: TabType,
	label: string,
	theme: ThemeKey,
}

type Language = {
	name: string,
	extension: string,
}

type TabType = 'statement' | 'table'
type WorkspaceScreen = 'statement' | 'submission_history'
type ThemeKey = 'pizza' | 'soc'
type WorkspaceView = {
	id: WorkspaceScreen,
	label: string,
}

export type {
	Task,
	Language,
	WorkspaceScreen,
	WorkspaceView,
	ThemeKey,
}
