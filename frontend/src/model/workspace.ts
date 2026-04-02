import {computed, withAsyncData, wrap,} from '@reatom/core'
import type {PlatformApi} from '../api/types'
import {defAction, defAtom, defComputed,} from '../common/createModelProvider'
import {verify} from '../common/exception'
import type {Task, WorkspaceScreen, WorkspaceView} from './types'

type DefWorkspaceModelArgs = {
    api: PlatformApi,
    tabs: Task[],
    views: WorkspaceView[],
}

function defWorkspaceModel({tabs, views, api}: DefWorkspaceModelArgs) {
    const tabsAtom = defAtom(tabs)
    const viewsAtom = defAtom(views)
    const selectedTaskIdAtom = defAtom(verify(tabs[0]).id)
    const selectedScreenAtom = defAtom<WorkspaceScreen>(verify(views[0]).id)
    const selectedTab = defComputed(() => {
        return tabsAtom().find(tab => tab.id === selectedTaskIdAtom())
    })
    const setSelectedTab = defAction(async (id: string) => {
        selectedTaskIdAtom.set(id)
        selectedScreenAtom.set(verify(viewsAtom()[0]).id)
    })
    const setSelectedScreen = defAction((screen: WorkspaceScreen) => {
        selectedScreenAtom.set(screen)
    })

    const dataAtom = computed(async () => {
        const id = selectedTaskIdAtom()
        return wrap(api.getData<string>({id}))
    }).extend(withAsyncData())

    return {
        dataAtom,
        tabsAtom,
        viewsAtom,
        selectedTaskIdAtom,
        selectedScreenAtom,
        setSelectedTab,
        setSelectedScreen,
        selectedTab,
    }
}

export {
    defWorkspaceModel,
}
