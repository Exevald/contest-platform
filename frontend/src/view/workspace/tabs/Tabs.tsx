import {reatomComponent} from '@reatom/react'
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
		<div className={joinStyles(className, styles.tabs)}>
			{tabsAtom().map(tab => {
				const selected = selectedTaskIdAtom() === tab.id
				return <div
					key={tab.id}
					onClick={() => setSelectedTab(tab.id)}
					className={joinStyles(
						styles.tab,
						selected && styles.selected,
					)}
				>
					{tab.label}
				</div>
			})}
		</div>
	)
})

export {
	Tabs,
}
