import {Badge} from '../../../../common/components/badge/Badge'
import {Card} from '../../../../common/components/card/Card'
import styles from './TaskPreview.module.css'

type TaskPreviewProps = {
    taskId: string,
}

function TaskPreview({taskId}: TaskPreviewProps) {
    switch (taskId) {
        case 'pizza-task-1':
            return <TaskOnePreview/>
        case 'pizza-task-2':
            return <TaskTwoPreview/>
        case 'pizza-task-3':
            return <TaskThreePreview/>
        default:
            return null
    }
}

function TaskOnePreview() {
    const rows = [
        ['6', 'Иван', 'active', 'Веселая', '10', '5', 'Пепперони', '2'],
        ['3', 'Данил', 'active', 'Зелёная', '10', '5', 'Грибная', '2'],
        ['1', 'Сима', 'active', 'Уродливая', '7', '1', 'Острая', '2'],
        ['1', 'Сима', 'active', 'Уродливая', '7', '1', 'Вкусная', '2'],
        ['5', 'Клиент 5', 'active', 'Космическая', '1', 'None', 'Пепперони', '3'],
    ]

    return (
        <div className={styles.previewStack}>
            <div className={styles.previewHero}>
                <div>
                    <div className={styles.previewEyebrow}>OUTPUT PREVIEW</div>
                    <div className={styles.previewTitle}>Clean Orders CSV</div>
                </div>
                <div className={styles.badgeRow}>
                    <Badge tone="ok">5 валидных строк</Badge>
                    <Badge tone="neutral">CSV → Table</Badge>
                </div>
            </div>

            <Card className={styles.tableShell}>
                <div className={styles.tableToolbar}>
                    <div className={styles.dotGroup}>
                        <span/>
                        <span/>
                        <span/>
                    </div>
                    <div className={styles.toolbarLabel}>rocket_slice.cleaned_orders.csv</div>
                </div>
                <div className={styles.scrollWrap}>
                    <table className={styles.table}>
                        <thead>
                        <tr>
                            <th>order_id</th>
                            <th>client_name</th>
                            <th>status</th>
                            <th>street</th>
                            <th>house</th>
                            <th>flat</th>
                            <th>item_name</th>
                            <th>count</th>
                        </tr>
                        </thead>
                        <tbody>
                        {rows.map(row => (
                            <tr key={row.join('-')}>
                                {row.map((cell, index) => (
                                    <td key={cell + index}>
                                        {index === 2 ? <Badge tone="ok">{cell}</Badge> : cell}
                                    </td>
                                ))}
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            </Card>
        </div>
    )
}

function TaskTwoPreview() {
    const regions = [
        {name: 'Восточный', revenue: '6600', width: '100%', tone: 'hot'},
        {name: 'Западный', revenue: '5500', width: '83%', tone: 'warm'},
        {name: 'Unknown', revenue: '2200', width: '33%', tone: 'cold'},
    ]

    return (
        <div className={styles.previewStack}>
            <div className={styles.previewHero}>
                <div>
                    <div className={styles.previewEyebrow}>OUTPUT PREVIEW</div>
                    <div className={styles.previewTitle}>Demand Heatmap + Report CSV</div>
                </div>
                <Badge tone="warning">Revenue sorted</Badge>
            </div>

            <div className={styles.regionGrid}>
                {regions.map(region => (
                    <Card key={region.name} className={styles.regionCard}>
                        <div className={styles.regionHeader}>
                            <span>{region.name}</span>
                            <span>{region.revenue}</span>
                        </div>
                        <div className={styles.heatTrack}>
                            <div
                                className={styles[`heatFill${region.tone === 'hot' ? 'Hot' : region.tone === 'warm' ? 'Warm' : 'Cold'}`]}
                                style={{width: region.width}}
                            />
                        </div>
                    </Card>
                ))}
            </div>

            <Card className={styles.tableShell}>
                <div className={styles.tableToolbar}>
                    <div className={styles.toolbarLabel}>region_report.csv</div>
                    <Badge tone="neutral">5 order rows</Badge>
                </div>
                <div className={styles.scrollWrap}>
                    <table className={styles.table}>
                        <thead>
                        <tr>
                            <th>region</th>
                            <th>orders_count</th>
                            <th>total_revenue</th>
                            <th>order_id</th>
                            <th>street</th>
                            <th>home</th>
                            <th>flat</th>
                            <th>price</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr>
                            <td>Восточный</td>
                            <td>2</td>
                            <td>6600</td>
                            <td>1</td>
                            <td>Уродливая</td>
                            <td>7</td>
                            <td>1</td>
                            <td>3400</td>
                        </tr>
                        <tr>
                            <td>Восточный</td>
                            <td>2</td>
                            <td>6600</td>
                            <td>1</td>
                            <td>Уродливая</td>
                            <td>7</td>
                            <td>1</td>
                            <td>3200</td>
                        </tr>
                        <tr>
                            <td>Западный</td>
                            <td>2</td>
                            <td>5500</td>
                            <td>3</td>
                            <td>Зелёная</td>
                            <td>10</td>
                            <td>5</td>
                            <td>2100</td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </Card>
        </div>
    )
}

function TaskThreePreview() {
    return (
        <div className={styles.previewStack}>
            <div className={styles.previewHero}>
                <div>
                    <div className={styles.previewEyebrow}>OUTPUT PREVIEW</div>
                    <div className={styles.previewTitle}>Courier Dispatch Board</div>
                </div>
                <Badge tone="ok">1 район выбран</Badge>
            </div>

            <div className={styles.dispatchGrid}>
                <Card className={styles.dispatchMain}>
                    <div className={styles.dispatchLabel}>Самый выгодный район</div>
                    <div className={styles.dispatchRegion}>Восточный</div>
                    <div className={styles.dispatchMeta}>Общая выручка: 6600</div>
                    <div className={styles.dispatchMeta}>Нужно курьеров: 1</div>
                </Card>
                <Card className={styles.dispatchSide}>
                    <div className={styles.dispatchLabel}>Вместимость смены</div>
                    <div className={styles.courierRow}>
                        {Array.from({length: 10}, (_, index) => (
                            <div
                                key={index}
                                className={index === 0 ? styles.courierActive : styles.courierIdle}
                            >
                                {index + 1}
                            </div>
                        ))}
                    </div>
                    <div className={styles.dispatchHint}>Один курьер берет до 3 заказов за рейс</div>
                </Card>
            </div>

            <Card className={styles.terminalCard}>
                <div className={styles.terminalLine}>Восточный</div>
                <div className={styles.terminalLine}>6600</div>
                <div className={styles.terminalLine}>1</div>
            </Card>
        </div>
    )
}

export {
    TaskPreview,
}
