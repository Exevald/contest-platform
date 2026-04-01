import {resolve} from 'path'
import react from '@vitejs/plugin-react'
import {defineConfig} from 'vite'

export default defineConfig(({command}) => ({
	plugins: [
		react(),
		{
			name: 'contestplatform-html-entry',
			transformIndexHtml(html) {
				const entryFile = command === 'serve'
					? '/src/index-dev.ts'
					: '/src/index.ts'

				return html.replace('/src/index-dev.ts', entryFile)
			},
		},
	],
	css: {
		modules: {
			localsConvention: 'camelCase',
		},
	},
	build: {
		outDir: resolve('dist'),
		emptyOutDir: true,
		target: 'es2022',
	},
}))
