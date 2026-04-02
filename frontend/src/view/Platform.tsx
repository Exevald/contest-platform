import {useState} from 'react'
import {getThemeTasks} from '../config/theme'
import {PlatformContext} from '../model/context'
import type {DefPlatformModel} from '../model/model'
import type {StartupData} from '../api/types'
import {PlatformView} from './PlatformView'
import {StartScreen} from './startScreen/StartScreen'

function Platform(props: DefPlatformModel) {
	const [startupData, setStartupData] = useState<StartupData>(props.startupData)
	const selectedTheme = startupData.selectedTheme

	if (!selectedTheme) {
		return (
			<StartScreen
				api={props.api}
				startupData={startupData}
				onStart={setStartupData}
			/>
		)
	}

	return (
		<PlatformContext
			api={props.api}
			startupData={{
				...startupData,
				tasks: getThemeTasks(selectedTheme as 'pizza' | 'soc').map(task => ({
					id: task.id,
					type: 'statement',
					label: task.label,
					theme: task.theme,
				})),
			}}
		>
			<PlatformView />
		</PlatformContext>
	)
}

export {
	Platform,
}
