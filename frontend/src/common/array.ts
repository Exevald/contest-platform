function insertBetween(arr: any[], separator: any) {
	if (arr.length <= 1) {
		return [...arr]
	}

	return arr.flatMap((item, index) =>
		index === arr.length - 1
			? [item]
			: [item, separator],
	)
}

export {
	insertBetween,
}
