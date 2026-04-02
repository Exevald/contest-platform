import type {HTMLAttributes, ReactNode} from 'react'
import {joinStyles} from '../../../common/joinStyles'
import styles from './Alert.module.css'

type AlertTone = 'neutral' | 'error'

type AlertProps = HTMLAttributes<HTMLDivElement> & {
	children: ReactNode,
	tone?: AlertTone,
}

function Alert({
	children,
	className,
	tone = 'neutral',
	...props
}: AlertProps) {
	return (
		<div
			{...props}
			className={joinStyles(styles.alert, styles[tone], className)}
		>
			{children}
		</div>
	)
}

export {
	Alert,
}
