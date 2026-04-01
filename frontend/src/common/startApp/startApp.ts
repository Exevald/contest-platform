import {createElement} from 'react'
import {createRoot} from 'react-dom/client'
import {PANIC} from '../exception'
import type {StartupApi} from './api/StartupApi'
import type {StartAppArgs} from './types'

const startApp = async <STARTUP_DATA, API extends StartupApi<STARTUP_DATA>>({
	api,
	appComponent,
	loaderComponent,
}: StartAppArgs<STARTUP_DATA, API>) => {
	const rootElement = document.getElementById('root')

	if (!rootElement) {
		PANIC('root element not found')
	}

	const root = createRoot(rootElement)

	if (loaderComponent) {
		root.render(createElement(loaderComponent))
	}

	try {
		const startupData = await api.getStartupData()

		root.render(
			createElement(appComponent, {
				api,
				startupData,
			}),
		)
	}
	catch (e) {
		PANIC(`Unable to load startup data or initialize localization. Error: ${JSON.stringify(e)}`)
	}
}

export {
	startApp,
}
