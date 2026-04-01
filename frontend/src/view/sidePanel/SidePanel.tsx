import {reatomComponent} from '@reatom/react'
import type {ChangeEvent} from 'react'
import type {ClassNameProps} from '../../common/components/types'
import {verify} from '../../common/exception'
import {joinStyles} from '../../common/joinStyles'
import {useModel} from '../../model/context'
import styles from './SidePanel.module.css'

const SidePanel = reatomComponent(({className}: ClassNameProps) => {
	const {
		selectedFile,
		selectedLanguageName,
		setSelectedFile,
		setSelectedLanguageName,
		handleSubmit,
		isSubmitDisabled,
		languagesAtom,
		titleAtom,
	} = useModel().sidePanel

	const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
		const file = event.target.files?.[0] || null
		setSelectedFile(file)
	}

	const handleLanguageChange = (event: ChangeEvent<HTMLSelectElement>) => {
		setSelectedLanguageName(event.target.value)
	}

	return (
		<div className={joinStyles(className, styles.container)}>
			<div className={styles.title}>{titleAtom()}</div>
			<input
				type="file"
				id="file-input"
				className={styles.fileInput}
				onChange={handleFileChange}
			/>
			<label htmlFor="file-input" className={styles.fileLabel}>
				<svg
					className={styles.fileIcon}
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
				<span className={styles.fileName}>
					{selectedFile()
						? verify(selectedFile()).name
						: 'Прикрепить файл'}
				</span>
			</label>

			<select
				className={styles.select}
				value={selectedLanguageName()}
				onChange={handleLanguageChange}
			>
				{languagesAtom().map(lang => (
					<option key={lang.name} value={lang.name}>
						{lang.name}
					</option>
				))}
			</select>

			<button
				className={styles.submitButton}
				onClick={handleSubmit}
				disabled={isSubmitDisabled()}
			>
					Отправить
			</button>
		</div>
	)
})

export {
	SidePanel,
}
