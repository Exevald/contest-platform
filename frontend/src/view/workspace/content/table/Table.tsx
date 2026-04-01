import {reatomComponent} from '@reatom/react'
import {useModel} from '../../../../model/context'
import {parseCsv} from '../../../../model/parse/parseCsv'
import {Preloader} from '../../preloader/Preloader'
import {TableView} from './TableView'

const Table = reatomComponent(() => {
	const {dataAtom} = useModel().workspace
	const data = parseCsv(dataAtom.data())
	const columns = Object.keys(data?.[0] || {}).map(key => ({
		key,
		title: key,
	}))
	return (
		dataAtom.ready()
			? <TableView data={data} columns={columns} />
			: <Preloader />
	)
})

export {
	Table,
}
