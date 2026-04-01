import type {Falsy} from '@reatom/core'
import {isTruthy} from './guards'

export function joinStyles(...classNames: (string | Falsy)[]): string {
	return classNames
		.filter(isTruthy)
		.join(' ')
}
