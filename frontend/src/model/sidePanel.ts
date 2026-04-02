import {
	computed, withAsyncData, wrap,
} from '@reatom/core'
import {
	defAction, defAtom, defComputed,
} from '../common/createModelProvider'
import {PANIC, verify} from '../common/exception'
import type {PlatformApi, SubmissionStatus} from '../api/types'
import type {Language} from './types'

type DefSidePanelModelArgs = {
	api: PlatformApi,
	getSelectedTaskId: () => string,
	title: string,
	participantCode: string,
	languages: Language[],
}

function defSidePanelModel({
	api,
	getSelectedTaskId,
	title,
	participantCode,
	languages,
}: DefSidePanelModelArgs) {
	if (languages.length === 0) {
		PANIC('No languages specified.')
	}
	const titleAtom = defAtom(title)
	const participantCodeAtom = defAtom(participantCode)

	const selectedFile = defAtom<File | null>(null)
	const selectedLanguageName = defAtom(verify(languages[0]).name)
	const submitErrorAtom = defAtom<string | null>(null)
	const submitResultAtom = defAtom<string | null>(null)
	const isSubmittingAtom = defAtom(false)
	const submissionRefreshAtom = defAtom(0)
	const expandedSubmissionIDsAtom = defAtom<string[]>([])
	let pollTimeout: ReturnType<typeof setTimeout> | null = null
	const setSelectedLanguageName = defAction((language: string) => {
		selectedLanguageName.set(language)
	})

	const languagesAtom = defAtom(languages)

	const setSelectedFile = defAction((file: File | null) => {
		selectedFile.set(file)
		languagesAtom().forEach(language => {
			if (file?.name.toLowerCase().endsWith(`.${language.extension.toLowerCase()}`)) {
				setSelectedLanguageName(language.name)
			}
		})
	})

	const submissionHistoryAtom = computed(async () => {
		submissionRefreshAtom()
		const id = getSelectedTaskId()
		return wrap(api.getSubmissionHistory<SubmissionStatus[]>({id}))
	}).extend(withAsyncData())

	const toggleSubmissionExpanded = defAction((submissionID: string) => {
		expandedSubmissionIDsAtom.set(value => (
			value.includes(submissionID)
				? value.filter(id => id !== submissionID)
				: [...value, submissionID]
		))
	})

	const scheduleStatusPoll = (submissionId: string) => {
		if (pollTimeout) {
			clearTimeout(pollTimeout)
		}

		const poll = async () => {
			try {
				const status = await api.getSubmissionStatus<SubmissionStatus | null>({id: submissionId})
				submissionRefreshAtom.set(value => value + 1)

				if (!status || !['PENDING', 'RUNNING', 'COMPILING'].includes(status.verdict)) {
					pollTimeout = null
					return
				}
			}
			catch {
				submissionRefreshAtom.set(value => value + 1)
			}

			pollTimeout = setTimeout(poll, 1000)
		}

		pollTimeout = setTimeout(poll, 300)
	}

	const handleSubmit = defAction(async () => {
		const file = selectedFile()
		const language = selectedLanguageName()

		if (!file) {
			return
		}

		isSubmittingAtom.set(true)
		submitErrorAtom.set(null)
		submitResultAtom.set(null)

		try {
			const text = await file.text()
			const result = await api.sendFile<{submission_id?: string, error?: string}>({
				id: getSelectedTaskId(),
				language,
				text,
			})

			if (result?.error) {
				submitErrorAtom.set(result.error)
				return
			}

			submitResultAtom.set(result?.submission_id ?? 'submitted')
			submissionRefreshAtom.set(value => value + 1)
			if (result?.submission_id) {
				scheduleStatusPoll(result.submission_id)
			}
		}
		catch (error) {
			submitErrorAtom.set(error instanceof Error ? error.message : String(error))
		}
		finally {
			isSubmittingAtom.set(false)
		}
	})

	const isSubmitDisabled = defComputed(() => !selectedFile() || isSubmittingAtom())

	return {
		titleAtom,
		participantCodeAtom,
		selectedFile,
		setSelectedLanguageName,
		setSelectedFile,
		selectedLanguageName,
		handleSubmit,
		isSubmitDisabled,
		languagesAtom,
		submitErrorAtom,
		submitResultAtom,
		isSubmittingAtom,
		submissionHistoryAtom,
		expandedSubmissionIDsAtom,
		toggleSubmissionExpanded,
	}
}

export {
	defSidePanelModel,
}
