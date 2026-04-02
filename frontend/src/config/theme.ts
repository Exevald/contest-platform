import type {Task} from '../model/types'

type ThemeKey = 'pizza' | 'soc'
type BadgeTone = 'ok' | 'error' | 'warning' | 'neutral'

type StatementTask = {
	id: string,
	theme: ThemeKey,
	label: string,
	status: string,
	statusTone: BadgeTone,
}

type ThemeOption = {
	id: ThemeKey,
	label: string,
}

const themeOptions: ThemeOption[] = [
	{id: 'pizza', label: 'RocketSlice'},
	{id: 'soc', label: 'Enterprise SOC'},
]

const tasks: StatementTask[] = [
	{id: 'pizza-task-1', theme: 'pizza', label: 'Фильтр заказов', status: 'CSV OUTPUT', statusTone: 'ok'},
	{id: 'pizza-task-2', theme: 'pizza', label: 'Карта спроса', status: 'REPORT OUTPUT', statusTone: 'warning'},
	{id: 'pizza-task-3', theme: 'pizza', label: 'Генератор маршрутов', status: 'DISPATCH OUTPUT', statusTone: 'ok'},
	{id: 'soc-task-1', theme: 'soc', label: 'Threat Pipeline', status: 'LOG STREAM', statusTone: 'error'},
	{id: 'soc-task-2', theme: 'soc', label: 'Threat Graph', status: 'SUSPICIOUS IP OUTPUT', statusTone: 'warning'},
	{id: 'soc-task-3', theme: 'soc', label: 'WAF Policies', status: 'MITIGATION OUTPUT', statusTone: 'ok'},
]

function getAllStartupTasks(): Task[] {
	return tasks.map(task => ({
		id: task.id,
		type: 'statement',
		label: task.label,
		theme: task.theme,
	}))
}

function getThemeOptions() {
	return themeOptions
}

function getTaskStatement(taskId: string) {
	return tasks.find(task => task.id === taskId)
}

function getThemeTasks(theme: ThemeKey) {
	return tasks.filter(task => task.theme === theme)
}

export {
	getAllStartupTasks,
	getTaskStatement,
	getThemeOptions,
	getThemeTasks,
}

export type {
	StatementTask,
	ThemeKey,
	ThemeOption,
}
