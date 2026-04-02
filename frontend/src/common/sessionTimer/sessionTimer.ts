import type {ThemeKey} from '../../model/types'

const TIMER_DURATION_MS = 3.5 * 60 * 60 * 1000

function getTimerStorageKey(participantCode: string, theme: ThemeKey) {
	return `contestplatform:timer:${participantCode}:${theme}`
}

function ensureTimerDeadline(participantCode: string, theme: ThemeKey) {
	if (!participantCode || !theme) {
		return 0
	}

	const key = getTimerStorageKey(participantCode, theme)
	const savedValue = window.localStorage.getItem(key)
	const parsedValue = savedValue ? Number(savedValue) : NaN

	if (Number.isFinite(parsedValue) && parsedValue > 0) {
		return parsedValue
	}

	const deadlineAt = Date.now() + TIMER_DURATION_MS
	window.localStorage.setItem(key, String(deadlineAt))
	return deadlineAt
}

export {
	ensureTimerDeadline,
	TIMER_DURATION_MS,
}
