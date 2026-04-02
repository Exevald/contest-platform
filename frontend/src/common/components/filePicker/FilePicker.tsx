import type {ChangeEvent} from 'react'
import styles from './FilePicker.module.css'

type FilePickerProps = {
    id: string,
    fileName: string,
    placeholder: string,
    onChange: (file: File | null) => void,
}

function FilePicker({
                        id,
                        fileName,
                        placeholder,
                        onChange,
                    }: FilePickerProps) {
    const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
        onChange(event.target.files?.[0] || null)
    }

    return (
        <div>
            <input
                type="file"
                id={id}
                className={styles.input}
                onChange={handleChange}
            />
            <label htmlFor={id} className={styles.trigger}>
                <svg
                    className={styles.icon}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
                    />
                </svg>
                <span className={styles.fileName}>{fileName || placeholder}</span>
            </label>
        </div>
    )
}

export {
    FilePicker,
}
