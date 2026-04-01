import {
	defAction, defAtom, defComputed,
} from '../common/createModelProvider'
import {PANIC, verify} from '../common/exception'
import type {Language} from './types'

type DefSidePanelModelArgs = {
	title: string,
	languages: Language[],
}

function defSidePanelModel({title, languages}: DefSidePanelModelArgs) {
	if (languages.length === 0) {
		PANIC('No languages specified.')
	}
	const titleAtom = defAtom(title)

	const selectedFile = defAtom<File | null>(null)
	const selectedLanguageName = defAtom(verify(languages[0]).name)
	const setSelectedLanguageName = defAction((language: string) => {
		selectedLanguageName.set(language)
	})

	const languagesAtom = defAtom(languages)

	const setSelectedFile = defAction((file: File | null) => {
		selectedFile.set(file)
		languagesAtom().forEach(language => {
			if (file?.name.endsWith(language.extension)) {
				setSelectedLanguageName(language.name)
			}
		})
	})

	const handleSubmit = defAction(() => {
		const file = selectedFile()
		const language = selectedLanguageName()

		if (file) {
			console.log('Отправка файла:', {
				file,
				language,
			})
		}
	})

	const isSubmitDisabled = defComputed(() => !selectedFile())

	return {
		titleAtom,
		selectedFile,
		setSelectedLanguageName,
		setSelectedFile,
		selectedLanguageName,
		handleSubmit,
		isSubmitDisabled,
		languagesAtom,
	}
}

export {
	defSidePanelModel,
}
