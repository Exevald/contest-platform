import {isNone} from './guards'
import type {None} from './guards'

function PANIC(message: string): never {
	throw new Error(message)
}

function verify<T>(
	value: T | None,
	message: string = 'verify failed',
): T {
	if (isNone(value)) {
		PANIC(message)
	}
	return value
}

export {
	PANIC,
	verify,
}
