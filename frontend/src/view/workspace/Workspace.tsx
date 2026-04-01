import {reatomComponent} from '@reatom/react'
import type {ClassNameProps} from '../../common/components/types'
import {joinStyles} from '../../common/joinStyles'
import {useModel} from '../../model/context'
import {Table} from './content/table/Table'
import {Tabs} from './tabs/Tabs'
import styles from './Workspace.module.css'

const Workspace = reatomComponent(({className}: ClassNameProps) => {
	const {selectedTab} = useModel().workspace

	let content

	switch (selectedTab()?.type) {
		case 'table':
			content = <Table />
			break
	}

	return (
		<div className={joinStyles(className, styles.workspace)}>
			<Tabs className={styles.tabs} />
			{content}
		</div>
	)
})

export {
	Workspace,
}
