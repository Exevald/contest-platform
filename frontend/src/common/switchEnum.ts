import {PANIC} from './exception'

function switchEnum<
	ENUM extends string,
	RESULT,
>(
	value: ENUM,
	handlers: Record<ENUM, RESULT>,
	getDefault?: (value: string) => RESULT,
): RESULT {
	if (!handlers.hasOwnProperty(value)) {
		if (getDefault) {
			return getDefault(value)
		}
		PANIC(`switchEnum() failed, received '${value}'`)
	}
	return handlers[value]
}

export {
	switchEnum,
}
