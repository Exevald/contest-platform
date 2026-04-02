import {
	defAtom, defComputed,
} from '../common/createModelProvider'

type DefTimerModelArgs = {
	deadlineAt: number,
}

function defTimerModel({deadlineAt}: DefTimerModelArgs) {
	const deadlineAtAtom = defAtom(deadlineAt)
	const nowAtom = defAtom(Date.now())

	setInterval(() => {
		nowAtom.set(Date.now())
	}, 1000)

	const remainingMsAtom = defComputed(() => {
		return Math.max(0, deadlineAtAtom() - nowAtom())
	})

	const isExpiredAtom = defComputed(() => remainingMsAtom() === 0)

	return {
		deadlineAtAtom,
		remainingMsAtom,
		isExpiredAtom,
	}
}

export {
	defTimerModel,
}
