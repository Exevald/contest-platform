import {useState} from 'react'
import {ensureTimerDeadline} from '../common/sessionTimer/sessionTimer'
import {getThemeTasks} from '../config/theme'
import {PlatformContext} from '../model/context'
import type {DefPlatformModel} from '../model/model'
import type {StartupData} from '../api/types'
import type {ThemeKey} from '../model/types'
import {PlatformView} from './PlatformView'
import {StartScreen} from './startScreen/StartScreen'

function withTimer(startupData: StartupData): StartupData {
    if (!startupData.participantCode || !startupData.selectedTheme) {
        return startupData
    }

    return {
        ...startupData,
        timerDeadlineAt: ensureTimerDeadline(
            startupData.participantCode,
            startupData.selectedTheme as ThemeKey,
        ),
    }
}

function Platform(props: DefPlatformModel) {
    const [startupData, setStartupData] = useState<StartupData>(() => withTimer(props.startupData))
    const selectedTheme = startupData.selectedTheme

    if (!selectedTheme) {
        return (
            <StartScreen
                api={props.api}
                startupData={startupData}
                onStart={data => setStartupData(withTimer(data))}
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
            <PlatformView/>
        </PlatformContext>
    )
}

export {
    Platform,
}
