import type {FC} from 'react'
import type {MessagesStartupData, StartupApi} from './api/StartupApi'

type AppComponentProps<API, STARTUP_DATA extends MessagesStartupData> = {
	api: API,
	startupData: STARTUP_DATA,
}

type StartAppArgs<STARTUP_DATA extends MessagesStartupData, API extends StartupApi<STARTUP_DATA>> = {
	api: API,
	appComponent: FC<AppComponentProps<API, STARTUP_DATA>>,
	loaderComponent?: FC,
}

export {
	AppComponentProps,
	StartAppArgs,
}
