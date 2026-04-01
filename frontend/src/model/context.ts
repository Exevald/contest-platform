import {createModelProvider} from '../common/createModelProvider'
import {defPlatformModel} from './model'
import type {DefPlatformModel, PlatformModel} from './model'

const {
	Provider: PlatformContext,
	useModel: usePlatformModel,
} = createModelProvider<DefPlatformModel, PlatformModel>(defPlatformModel)

const useModel = usePlatformModel

export {
	PlatformContext,
	useModel,
}
