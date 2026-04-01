import {reatomComponent} from '@reatom/react'
import {joinStyles} from '../../joinStyles'
import styles from './Divider.module.css'

type DividerProps = {
	type?: 'horizontal' | 'vertical',
}

const Divider = reatomComponent(({type = 'vertical'}: DividerProps) => {
	return (
		<div className={joinStyles(styles.divider, styles[type])} />
	)
})

export {
	Divider,
}
