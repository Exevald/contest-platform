import {reatomComponent} from '@reatom/react'
import {Alert} from '../../common/components/alert/Alert'
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
        participantCodeAtom,
        isTimerExpiredAtom,
        submitErrorAtom,
        submitResultAtom,
        isSubmittingAtom,
    } = useModel().sidePanel

    return (
        <div className={joinStyles(className, styles.container)}>
            <div className={styles.title}>{titleAtom()}</div>
            <div className={styles.participantCard}>
                <div className={styles.participantLabel}>Код участника</div>
                <div className={styles.participantValue}>{participantCodeAtom()}</div>
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

            {isTimerExpiredAtom() ? (
                <Alert tone="error">Время вышло. Отправка решений недоступна.</Alert>
            ) : null}

            {submitErrorAtom() ? (
                <Alert tone="error">{submitErrorAtom()}</Alert>
            ) : null}
            {submitResultAtom() ? (
                <Alert>Посылка: {submitResultAtom()}</Alert>
            ) : null}
        </div>
    )
})

export {
    SidePanel,
}
