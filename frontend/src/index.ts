import {api} from './api/cef/api'
import {startApp} from './common/startApp/startApp'
import {Platform} from './view/Platform'

await startApp({
	appComponent: Platform,
	api,
})
