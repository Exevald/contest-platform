import type {ReactNode} from 'react'
import styles from './TableView.module.css'

type Column<T> = {
    key: keyof T | string,
    title: string,
    render?: (value: T[keyof T], record: T) => ReactNode,
}

type TableViewProps<T> = {
    columns: Column<T>[],
    data: T[],
    rowKey?: keyof T | string,
}

function TableView<T extends Record<string, unknown>>({
                                                          columns,
                                                          data,
                                                          rowKey = 'id',
                                                      }: TableViewProps<T>) {
    return (
        <div className={styles.tableContainer}>
            <table className={styles.table}>
                <thead className={styles.head}>
                <tr className={styles.row}>
                    {columns.map(column => (
                        <th key={String(column.key)} className={styles.headerCell}>
                            {column.title}
                        </th>
                    ))}
                </tr>
                </thead>
                <tbody>
                {data.map((record, index) => (
                    <tr
                        key={String((record as any)[rowKey] ?? index)}
                        className={styles.row}
                    >
                        {columns.map(column => {
                            const value = (record as any)[column.key]
                            const content = column.render
                                ? column.render(value, record)
                                : String(value ?? '')

                            return (
                                <td key={String(column.key)} className={styles.cell}>
                                    {content}
                                </td>
                            )
                        })}
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    )
}

export {
    TableView,
}
export type {
    Column,
    TableViewProps,
}
