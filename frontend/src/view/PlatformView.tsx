import {reatomComponent} from '@reatom/react'
import {Divider} from '../common/components/divider/Divider'
import styles from './PlatformView.module.css'
import {SidePanel} from './sidePanel/SidePanel'
import {TimerBanner} from './timerBanner/TimerBanner'
import {Workspace} from './workspace/Workspace'

const PlatformView = reatomComponent(() => (
    <div className={styles.container}>
        <TimerBanner/>
        <SidePanel className={styles.sidePanel}/>
        <Divider/>
        <Workspace className={styles.workspace}/>
    </div>
))

export {
    PlatformView,
}
