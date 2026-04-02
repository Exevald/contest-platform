import {reatomComponent} from '@reatom/react'
import {Badge} from '../../../../common/components/badge/Badge'
import {Button} from '../../../../common/components/button/Button'
import {Card} from '../../../../common/components/card/Card'
import {useModel} from '../../../../model/context'
import {Preloader} from '../../preloader/Preloader'
import styles from './SubmissionHistory.module.css'

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

function verdictTone(verdict: string) {
    switch (verdict) {
        case 'OK':
            return 'ok'
        case 'TLE':
        case 'MLE':
            return 'warning'
        case 'WA':
        case 'RE':
        case 'CE':
        case 'SE':
            return 'error'
        default:
            return 'neutral'
    }
}

const SubmissionHistory = reatomComponent(() => {
    const {
        submissionHistoryAtom,
        expandedSubmissionIDsAtom,
        toggleSubmissionExpanded,
    } = useModel().sidePanel

    if (!submissionHistoryAtom.ready()) {
        return <Preloader fullScreen={false}/>
    }

    const submissions = submissionHistoryAtom.data()
    const expandedSubmissionIDs = expandedSubmissionIDsAtom()

    if (submissions?.length === 0) {
        return <div className={styles.empty}>Ещё не было посылок</div>
    }

    return (
        <div className={styles.list}>
            {submissions?.map(submission => {
                const isExpanded = expandedSubmissionIDs.includes(submission.submissionId)

                return (
                    <Card key={submission.submissionId} className={styles.item}>
                        <Button
                            variant="ghost"
                            className={styles.header}
                            onClick={() => toggleSubmissionExpanded(submission.submissionId)}
                        >
                            <div className={styles.summary}>
                                <span className={styles.id}>{submission.submissionId}</span>
                                <span className={styles.meta}>
									{submission.language}
                                    {' • '}
                                    {new Date(submission.createdAt).toLocaleString()}
								</span>
                            </div>
                            <div className={styles.headerRight}>
                                <Badge tone={verdictTone(submission.verdict)}>
                                    {verdictLabelMap[submission.verdict] ?? submission.verdict}
                                </Badge>
                                <span className={styles.chevron}>{isExpanded ? '−' : '+'}</span>
                            </div>
                        </Button>
                        {isExpanded ? (
                            <div className={styles.body}>
                                <div className={styles.sectionTitle}>Результат компиляции</div>
                                <pre className={styles.output}>
									{submission.compilationOutput || 'Вывод компиляции пуст.'}
								</pre>
                            </div>
                        ) : null}
                    </Card>
                )
            })}
        </div>
    )
})

export {
    SubmissionHistory,
}
