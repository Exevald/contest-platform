import {reatomComponent} from '@reatom/react'
import {useModel} from '../../model/context'
import styles from './TimerBanner.module.css'

function formatRemaining(ms: number) {
    const totalSeconds = Math.ceil(ms / 1000)
    const hours = Math.floor(totalSeconds / 3600)
    const minutes = Math.floor((totalSeconds % 3600) / 60)
    const seconds = totalSeconds % 60

    return [hours, minutes, seconds]
        .map(value => String(value).padStart(2, '0'))
        .join(':')
}

const TimerBanner = reatomComponent(() => {
    const {remainingMsAtom, isExpiredAtom} = useModel().timer
    const remainingMs = remainingMsAtom()

    return (
        <div className={styles.wrap}>
            <div className={isExpiredAtom() ? styles.timerExpired : styles.timer}>
                <span className={styles.label}>Осталось времени</span>
                <strong className={styles.value}>
                    {formatRemaining(remainingMs)}
                </strong>
            </div>
        </div>
    )
})

export {
    TimerBanner,
}
