import {reatomComponent} from '@reatom/react'
import {Badge} from '../../../../common/components/badge/Badge'
import {Card} from '../../../../common/components/card/Card'
import {getTaskStatement} from '../../../../config/theme'
import {useModel} from '../../../../model/context'
import {SocTaskPreview} from './SocTaskPreview'
import styles from './Statement.module.css'
import {TaskPreview} from './TaskPreview'

const Statement = reatomComponent(() => {
	const {selectedTaskIdAtom} = useModel().workspace
	const task = getTaskStatement(selectedTaskIdAtom())

	if (!task) {
		return (
			<div className={styles.emptyState}>
				<Card className={styles.emptyCard}>
					<div className={styles.emptyTitle}>Задача не найдена</div>
					<div className={styles.emptyText}>
						Для выбранной темы не найден output preview.
					</div>
				</Card>
			</div>
		)
	}

	return (
		<div className={styles.statement}>
			<div className={styles.header}>
				<h1 className={styles.title}>{task.label}</h1>
				<Badge tone={task.statusTone}>{task.status}</Badge>
			</div>

			<Card className={styles.previewCard}>
				{task.theme === 'soc'
					? <SocTaskPreview taskId={task.id} />
					: <TaskPreview taskId={task.id} />}
			</Card>
		</div>
	)
})

export {
	Statement,
}
