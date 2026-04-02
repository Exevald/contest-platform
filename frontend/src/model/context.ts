import {createModelProvider} from '../common/createModelProvider'
import type {DefPlatformModel, PlatformModel} from './model'
import {defPlatformModel} from './model'

const {
    Provider: PlatformContext,
    useModel: usePlatformModel,
} = createModelProvider<DefPlatformModel, PlatformModel>(defPlatformModel)

const useModel = usePlatformModel

export {
    PlatformContext,
    useModel,
}
