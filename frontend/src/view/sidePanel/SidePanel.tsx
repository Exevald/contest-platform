import {reatomComponent} from '@reatom/react'
import {Button} from '../../common/components/button/Button'
import {FilePicker} from '../../common/components/filePicker/FilePicker'
import {Select} from '../../common/components/select/Select'
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
		submitErrorAtom,
		submitResultAtom,
		isSubmittingAtom,
	} = useModel().sidePanel

	return (
		<div className={joinStyles(className, styles.container)}>
			<div className={styles.title}>{titleAtom()}</div>
			<div className={styles.caption}>
				Загрузите C++ файл с решением. История посылок и результаты компиляции доступны справа.
			</div>

			<FilePicker
				id="file-input"
				fileName={selectedFile() ? verify(selectedFile()).name : ''}
				placeholder="Прикрепить файл"
				onChange={setSelectedFile}
			/>

			<Select
				value={selectedLanguageName()}
				options={languagesAtom().map(lang => ({
					value: lang.name,
					label: lang.name,
				}))}
				onChange={event => setSelectedLanguageName(event.target.value)}
			/>

			<Button
				block
				onClick={handleSubmit}
				disabled={isSubmitDisabled()}
			>
				{isSubmittingAtom() ? 'Отправка...' : 'Отправить'}
			</Button>

			{submitErrorAtom() ? (
				<div className={styles.messageError}>{submitErrorAtom()}</div>
			) : null}
			{submitResultAtom() ? (
				<div className={styles.messageSuccess}>Посылка: {submitResultAtom()}</div>
			) : null}
		</div>
	)
})

export {
	SidePanel,
}
