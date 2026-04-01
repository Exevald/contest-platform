/* eslint-disable spaced-comment */

declare module '*.module.css' {
	const classes: Record<string, string>
	export default classes
}

declare module '*.jpg' {
	const value: string
	export default value
}

declare module '*.png' {
	const value: string
	export default value
}

declare module '*.mp3' {
	const value: string
	export default value
}

declare module '*.css' {
	const content: Record<string, string>
	export default content
}
