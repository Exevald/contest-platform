import type {HTMLAttributes, ReactNode} from 'react'
import {joinStyles} from '../../joinStyles.ts'
import styles from './Card.module.css'

type CardProps = HTMLAttributes<HTMLDivElement> & {
    children: ReactNode,
}

function Card({
                  children,
                  className,
                  ...props
              }: CardProps) {
    return (
        <div {...props} className={joinStyles(styles.card, className)}>
            {children}
        </div>
    )
}

export {
    Card,
}
