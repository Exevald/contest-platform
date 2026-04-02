import type {HTMLAttributes, ReactNode} from 'react'
import {joinStyles} from '../../../common/joinStyles'
import styles from './Badge.module.css'

type BadgeTone = 'ok' | 'error' | 'warning' | 'neutral'

type BadgeProps = HTMLAttributes<HTMLSpanElement> & {
	children: ReactNode,
	tone?: BadgeTone,
}

function Badge({
	children,
	className,
	tone = 'neutral',
	...props
}: BadgeProps) {
	return (
		<span
			{...props}
			className={joinStyles(styles.badge, styles[tone], className)}
		>
			{children}
		</span>
	)
}

export {
	Badge,
}
