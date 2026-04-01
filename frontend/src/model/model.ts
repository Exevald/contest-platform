import type {PlatformApi, StartupData} from '../api/types'
import {defSidePanelModel} from './sidePanel'
import {defWorkspaceModel} from './workspace'

type DefPlatformModel = {
	api: PlatformApi,
	startupData: StartupData,
}

function defPlatformModel(args: DefPlatformModel) {
	return {
		sidePanel: defSidePanelModel(args.startupData),
		workspace: defWorkspaceModel({
			tabs: args.startupData.tasks,
			api: args.api,
		}),
	}
}

type PlatformModel = ReturnType<typeof defPlatformModel>

export {
	DefPlatformModel,
	PlatformModel,
	defPlatformModel,
}
