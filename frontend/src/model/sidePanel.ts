import {
	defAction, defAtom, defComputed,
} from '../common/createModelProvider'
import {PANIC, verify} from '../common/exception'

type DefSidePanelModelArgs = {
	title: string,
	languages: string[],
}

function defSidePanelModel({title, languages}: DefSidePanelModelArgs) {
	if (languages.length === 0) {
		PANIC('No languages specified.')
	}
	const titleAtom = defAtom(title)

	const selectedFile = defAtom<File | null>(null)
	const selectedLanguage = defAtom(verify(languages[0]))

	const languagesAtom = defAtom(languages)

	const setSelectedFile = defAction((file: File | null) => {
		selectedFile.set(file)
	})

	const setSelectedLanguage = defAction((language: string) => {
		selectedLanguage.set(language)
	})

	const handleSubmit = defAction(() => {
		const file = selectedFile()
		const language = selectedLanguage()

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
		selectedLanguage,
		setSelectedFile,
		setSelectedLanguage,
		handleSubmit,
		isSubmitDisabled,
		languagesAtom,
	}
}

export {
	defSidePanelModel,
}
