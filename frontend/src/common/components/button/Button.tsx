import type {ButtonHTMLAttributes, ReactNode} from 'react'
import {joinStyles} from '../../joinStyles.ts'
import styles from './Button.module.css'

type ButtonVariant = 'primary' | 'secondary' | 'ghost'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    children: ReactNode,
    variant?: ButtonVariant,
    block?: boolean,
}

function Button({
                    children,
                    className,
                    variant = 'primary',
                    block = false,
                    type = 'button',
                    ...props
                }: ButtonProps) {
    return (
        <button
            {...props}
            type={type}
            className={joinStyles(styles.button, styles[variant], block && styles.block, className)}
        >
            {children}
        </button>
    )
}

export {
    Button,
}
