import {reatomComponent} from '@reatom/react'
import type {ChangeEvent} from 'react'
import type {ClassNameProps} from '../../common/components/types'
import {verify} from '../../common/exception'
import {joinStyles} from '../../common/joinStyles'
import {useModel} from '../../model/context'
import styles from './SidePanel.module.css'

const verdictLabelMap: Record<string, string> = {
	CE: 'Compilation Error',
	MLE: 'Memory Limit',
	OK: 'Accepted',
	PENDING: 'Pending',
	RE: 'Runtime Error',
	RUNNING: 'Running',
	SE: 'System Error',
	TLE: 'Time Limit',
	WA: 'Wrong Answer',
}

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
		submitErrorAtom,
		submitResultAtom,
		isSubmittingAtom,
		latestSubmissionAtom,
	} = useModel().sidePanel

	const latestSubmission = latestSubmissionAtom.ready()
		? latestSubmissionAtom.data()
		: null

	const verdictClassName = latestSubmission
		? styles[`verdict${latestSubmission.verdict}` as keyof typeof styles] ?? styles.verdictDefault
		: styles.verdictDefault

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
				{isSubmittingAtom() ? 'Отправка...' : 'Отправить'}
			</button>
			{submitErrorAtom() ? (
				<div className={styles.fileName}>{submitErrorAtom()}</div>
			) : null}
			{submitResultAtom() ? (
				<div className={styles.fileName}>Submission: {submitResultAtom()}</div>
			) : null}
			<div className={styles.statusCard}>
				<div className={styles.statusTitle}>Последняя посылка</div>
				{latestSubmission ? (
					<>
						<div className={styles.statusRow}>
							<span className={styles.statusLabel}>ID</span>
							<span className={styles.statusValue}>{latestSubmission.submissionId}</span>
						</div>
						<div className={styles.statusRow}>
							<span className={styles.statusLabel}>Вердикт</span>
							<span className={joinStyles(styles.verdictBadge, verdictClassName)}>
								{verdictLabelMap[latestSubmission.verdict] ?? latestSubmission.verdict}
							</span>
						</div>
						<div className={styles.statusRow}>
							<span className={styles.statusLabel}>Тесты</span>
							<span className={styles.statusValue}>
								{latestSubmission.testsPassed}/{latestSubmission.testsTotal}
							</span>
						</div>
						<div className={styles.statusRow}>
							<span className={styles.statusLabel}>Язык</span>
							<span className={styles.statusValue}>{latestSubmission.language}</span>
						</div>
					</>
				) : (
					<div className={styles.statusEmpty}>Ещё не было посылок</div>
				)}
			</div>
		</div>
	)
})

export {
	SidePanel,
}
