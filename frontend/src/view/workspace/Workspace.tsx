import {reatomComponent} from '@reatom/react'
import {Button} from '../../common/components/button/Button'
import type {ClassNameProps} from '../../common/components/types'
import {joinStyles} from '../../common/joinStyles'
import {useModel} from '../../model/context'
import {SubmissionHistory} from './content/submissionHistory/SubmissionHistory'
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
			content = <SubmissionHistory />
			break
		case 'statement':
		default:
			switch (selectedTab()?.type) {
				case 'table':
					content = <Table />
					break
			}
			break
	}

	return (
		<div className={joinStyles(className, styles.workspace)}>
			<Tabs className={styles.tabs} />
			<div className={styles.screenTabs}>
				{viewsAtom().map(view => (
					<Button
						key={view.id}
						variant={selectedScreenAtom() === view.id ? 'primary' : 'secondary'}
						className={joinStyles(styles.screenTab, selectedScreenAtom() === view.id && styles.screenTabSelected)}
						onClick={() => setSelectedScreen(view.id)}
					>
						{view.label}
					</Button>
				))}
			</div>
			<div className={styles.content}>{content}</div>
		</div>
	)
})

export {
	Workspace,
}
