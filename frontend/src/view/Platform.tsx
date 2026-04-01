import {PlatformContext} from '../model/context'
import type {DefPlatformModel} from '../model/model'
import {PlatformView} from './PlatformView'

function Platform(props: DefPlatformModel) {
	return (
		<PlatformContext {...props}>
			<PlatformView />
		</PlatformContext>
	)
}

export {
	Platform,
}
