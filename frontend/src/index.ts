import './index.css'
import {api} from './api'
import {startApp} from './common/startApp/startApp'
import {Platform} from './view/Platform'

await startApp({
    appComponent: Platform,
    api,
})
