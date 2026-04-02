import {useState} from 'react'
import {Button} from '../../common/components/button/Button'
import {SegmentedTabs} from '../../common/components/segmentedTabs/SegmentedTabs'
import type {PlatformApi, StartupData} from '../../api/types'
import {getThemeOptions} from '../../config/theme'
import type {ThemeKey} from '../../model/types'
import styles from './StartScreen.module.css'

type StartScreenProps = {
    api: PlatformApi,
    startupData: StartupData,
    onStart: (data: StartupData) => void,
}

function StartScreen({
                         api,
                         startupData,
                         onStart,
                     }: StartScreenProps) {
    const [participantCode, setParticipantCode] = useState(startupData.participantCode)
    const [theme, setTheme] = useState<ThemeKey>((startupData.selectedTheme as ThemeKey) || 'pizza')
    const [error, setError] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)

    const handleStart = async () => {
        if (!participantCode.trim()) {
            setError('Укажите код участника')
            return
        }

        setError('')
        setIsSubmitting(true)

        try {
            const nextData = await api.startSession({
                participantCode: participantCode.trim(),
                theme,
            })
            onStart(nextData)
        } catch (e) {
            setError(e instanceof Error ? e.message : 'Не удалось сохранить стартовую сессию')
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <div className={styles.page}>
            <div className={styles.panel}>
                <div className={styles.eyebrow}>ContestPlatform</div>
                <h1 className={styles.title}>Старт участника</h1>
                <p className={styles.subtitle}>
                    Сначала укажи код участника и тему. После старта откроется только выбранный пакет заданий.
                </p>

                <div className={styles.block}>
                    <div className={styles.label}>Код участника</div>
                    <input
                        value={participantCode}
                        onChange={event => setParticipantCode(event.target.value)}
                        className={styles.input}
                        placeholder="Например: A17-ROCKET"
                    />
                </div>

                <div className={styles.block}>
                    <div className={styles.label}>Тема</div>
                    <SegmentedTabs
                        className={styles.themeTabs}
                        itemClassName={styles.themeTab}
                        items={getThemeOptions()}
                        selectedId={theme}
                        onSelect={id => setTheme(id as ThemeKey)}
                    />
                </div>

                {error ? <div className={styles.error}>{error}</div> : null}

                <Button
                    block
                    className={styles.startButton}
                    onClick={handleStart}
                    disabled={isSubmitting}
                >
                    {isSubmitting ? 'Сохранение...' : 'Начать'}
                </Button>
            </div>
        </div>
    )
}

export {
    StartScreen,
}
