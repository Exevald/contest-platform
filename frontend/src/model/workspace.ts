import {
	computed, withAsyncData, wrap,
} from '@reatom/core'
import type {PlatformApi} from '../api/types'
import {
	defAction, defAtom, defComputed,
} from '../common/createModelProvider'
import {verify} from '../common/exception'
import type {Task} from './types'

type DefWorkspaceModelArgs = {
	api: PlatformApi,
	tabs: Task[],
}

function defWorkspaceModel({tabs, api}: DefWorkspaceModelArgs) {
	const tabsAtom = defAtom(tabs)
	const selectedTaskIdAtom = defAtom(verify(tabs[0]).id)
	const selectedTab = defComputed(() => {
		return tabsAtom().find(tab => tab.id === selectedTaskIdAtom())
	})
	const setSelectedTab = defAction(async (id: string) => {
		selectedTaskIdAtom.set(id)
	})

	const dataAtom = computed(async () => {
		const id = selectedTaskIdAtom()
		return wrap(api.getData<string>({id}))
	}).extend(withAsyncData())

	return {
		dataAtom,
		tabsAtom,
		selectedTaskIdAtom,
		setSelectedTab,
		selectedTab,
	}
}

export {
	defWorkspaceModel,
}
