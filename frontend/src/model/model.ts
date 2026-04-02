import type {PlatformApi, StartupData} from '../api/types'
import {defSidePanelModel} from './sidePanel'
import {defTimerModel} from './timer'
import {defWorkspaceModel} from './workspace'

type DefPlatformModel = {
	api: PlatformApi,
	startupData: StartupData,
}

function defPlatformModel(args: DefPlatformModel) {
	const timer = defTimerModel({
		deadlineAt: args.startupData.timerDeadlineAt ?? 0,
	})
	const workspace = defWorkspaceModel({
		tabs: args.startupData.tasks,
		views: args.startupData.workspaceViews,
		api: args.api,
	})

	return {
		sidePanel: defSidePanelModel({
			title: args.startupData.title,
			participantCode: args.startupData.participantCode,
			languages: args.startupData.languages,
			api: args.api,
			getSelectedTaskId: () => workspace.selectedTaskIdAtom(),
			isTimerExpired: () => timer.isExpiredAtom(),
		}),
		timer,
		workspace,
	}
}

type PlatformModel = ReturnType<typeof defPlatformModel>

export {
	DefPlatformModel,
	PlatformModel,
	defPlatformModel,
}
