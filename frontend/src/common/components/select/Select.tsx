import type {SelectHTMLAttributes} from 'react'
import {joinStyles} from '../../joinStyles.ts'
import styles from './Select.module.css'

type Option = {
    value: string,
    label: string,
}

type SelectProps = Omit<SelectHTMLAttributes<HTMLSelectElement>, 'children'> & {
    options: Option[],
}

function Select({
                    className,
                    options,
                    ...props
                }: SelectProps) {
    return (
        <div className={joinStyles(styles.wrapper, className)}>
            <select {...props} className={styles.select}>
                {options.map(option => (
                    <option key={option.value} value={option.value}>
                        {option.label}
                    </option>
                ))}
            </select>
            <span className={styles.chevron}>⌄</span>
        </div>
    )
}

export {
    Select,
}
