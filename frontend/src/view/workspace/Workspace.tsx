import {reatomComponent} from '@reatom/react'
import {SegmentedTabs} from '../../common/components/segmentedTabs/SegmentedTabs'
import type {ClassNameProps} from '../../common/components/types'
import {joinStyles} from '../../common/joinStyles'
import {useModel} from '../../model/context'
import {SubmissionHistory} from './content/submissionHistory/SubmissionHistory'
import {Statement} from './content/statement/Statement'
import {Table} from './content/table/Table'
import {Tabs} from './tabs/Tabs'
import styles from './Workspace.module.css'

const Workspace = reatomComponent(({className}: ClassNameProps) => {
    const {
        selectedTab,
        viewsAtom,
        selectedScreenAtom,
        setSelectedScreen,
    } = useModel().workspace

    let content

    switch (selectedScreenAtom()) {
        case 'submission_history':
            content = <SubmissionHistory/>
            break
        case 'statement':
        default:
            switch (selectedTab()?.type) {
                case 'statement':
                    content = <Statement/>
                    break
                case 'table':
                    content = <Table/>
                    break
            }
            break
    }

    return (
        <div className={joinStyles(className, styles.workspace)}>
            <Tabs className={styles.tabs}/>
            <SegmentedTabs
                className={styles.screenTabs}
                itemClassName={styles.screenTab}
                items={viewsAtom()}
                selectedId={selectedScreenAtom()}
                onSelect={setSelectedScreen}
            />
            <div className={styles.content}>{content}</div>
        </div>
    )
})

export {
    Workspace,
}
