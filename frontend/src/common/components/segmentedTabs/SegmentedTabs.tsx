import {Button} from '../button/Button'
import {joinStyles} from '../../../common/joinStyles'
import styles from './SegmentedTabs.module.css'

type SegmentedTab = {
	id: string,
	label: string,
}

type SegmentedTabsProps = {
	items: SegmentedTab[],
	selectedId: string,
	onSelect: (id: string) => void,
	className?: string,
	itemClassName?: string,
}

function SegmentedTabs({
	items,
	selectedId,
	onSelect,
	className,
	itemClassName,
}: SegmentedTabsProps) {
	return (
		<div className={joinStyles(styles.tabs, className)}>
			{items.map(item => (
				<Button
					key={item.id}
					variant={selectedId === item.id ? 'primary' : 'secondary'}
					className={joinStyles(styles.item, itemClassName)}
					onClick={() => onSelect(item.id)}
				>
					{item.label}
				</Button>
			))}
		</div>
	)
}

export {
	SegmentedTabs,
}
