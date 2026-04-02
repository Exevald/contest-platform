import {reatomComponent} from '@reatom/react'
import {SegmentedTabs} from '../../../common/components/segmentedTabs/SegmentedTabs'
import type {ClassNameProps} from '../../../common/components/types'
import {joinStyles} from '../../../common/joinStyles'
import {useModel} from '../../../model/context'
import styles from './Tabs.module.css'

const Tabs = reatomComponent(({className}: ClassNameProps) => {
    const {
        tabsAtom,
        selectedTaskIdAtom,
        setSelectedTab,
    } = useModel().workspace

    return (
        <SegmentedTabs
            className={joinStyles(className, styles.tabs)}
            itemClassName={styles.tab}
            items={tabsAtom().map(tab => ({
                id: tab.id,
                label: tab.label,
            }))}
            selectedId={selectedTaskIdAtom()}
            onSelect={setSelectedTab}
        />
    )
})

export {
    Tabs,
}
