type CsvOptions = {
	delimiter?: string,
	hasHeader?: boolean,
	trim?: boolean,
	skipEmptyLines?: boolean,
}

type CsvRow = {
	[key: string]: string,
}

type CsvResult = CsvRow[] | string[][]

function parseCsv(str: string, options: CsvOptions = {}): CsvResult {
	const {
		delimiter = ',',
		hasHeader = true,
		trim = true,
		skipEmptyLines = true,
	} = options

	if (!str || typeof str !== 'string') {
		return hasHeader
			? []
			: [[]]
	}

	const rows: string[][] = []
	let currentRow: string[] = []
	let currentField = ''
	let inQuotes = false
	let i = 0

	while (i < str.length) {
		const char = str[i]
		const nextChar = str[i + 1]

		if (inQuotes) {
			if (char === '"') {
				if (nextChar === '"') {
					// Экранированная кавычка
					currentField += '"'
					i += 2
					continue
				}
				else {
					// Конец кавычек
					inQuotes = false
					i++
					continue
				}
			}
			else {
				currentField += char
				i++
				continue
			}
		}
		else {
			if (char === '"') {
				inQuotes = true
				i++
				continue
			}

			if (char === delimiter) {
				currentRow.push(trim
					? currentField.trim()
					: currentField)
				currentField = ''
				i++
				continue
			}

			if (char === '\r' && nextChar === '\n') {
				currentRow.push(trim
					? currentField.trim()
					: currentField)
				currentField = ''
				if (currentRow.length > 0 && (!skipEmptyLines || currentRow.some(f => f !== ''))) {
					rows.push(currentRow)
				}
				currentRow = []
				i += 2
				continue
			}

			if (char === '\n' || char === '\r') {
				currentRow.push(trim
					? currentField.trim()
					: currentField)
				currentField = ''
				if (currentRow.length > 0 && (!skipEmptyLines || currentRow.some(f => f !== ''))) {
					rows.push(currentRow)
				}
				currentRow = []
				i++
				continue
			}

			currentField += char
			i++
			continue
		}
	}

	// Добавляем последнее поле и строку
	if (currentField || currentRow.length > 0) {
		currentRow.push(trim
			? currentField.trim()
			: currentField)
		if (currentRow.length > 0 && (!skipEmptyLines || currentRow.some(f => f !== ''))) {
			rows.push(currentRow)
		}
	}

	if (!hasHeader || rows.length === 0) {
		return rows as string[][]
	}

	const headers = rows[0]
	const dataRows = rows.slice(1)

	return dataRows.map(row => {
		const obj: CsvRow = {}
		headers.forEach((header, index) => {
			obj[header] = row[index] ?? ''
		})
		return obj
	})
}

export {
	parseCsv,
}
