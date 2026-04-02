import {reatomComponent} from '@reatom/react'
import {Button} from '../../../common/components/button/Button'
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
				return <Button
					key={tab.id}
					onClick={() => setSelectedTab(tab.id)}
					variant={selected ? 'primary' : 'secondary'}
					className={joinStyles(
						styles.tab,
						selected && styles.selected,
					)}
				>
					{tab.label}
				</Button>
			})}
		</div>
	)
})

export {
	Tabs,
}
