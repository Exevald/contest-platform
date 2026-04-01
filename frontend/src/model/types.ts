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

export type {
	Task,
	Language,
}
