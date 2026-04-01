type None = null | undefined
type Falsy = false | None

function isUndef<T>(value: T | undefined): value is undefined {
	return value === undefined
}

function isNull<T>(value: T | null): value is null {
	return value === null
}

function isNone<T>(value: T | None): value is None {
	return isUndef(value) || isNull(value)
}

function isTruthy<T>(value: T | Falsy): value is T {
	return !isFalsy(value)
}

function isFalsy<T>(value: T | Falsy): value is Falsy {
	return value === false || isNone(value)
}

export {
	None,
	isNone,
	isNull,
	isUndef,
	isFalsy,
	isTruthy,
}
