import type {Plugin} from 'vite'

function generateIndexHtml(): Plugin {
	return {
		name: 'generate-index-html',
		apply: 'build',
		enforce: 'post',
		generateBundle(_, bundle) {
			const cssFiles = Object.keys(bundle).filter(
				file => file.endsWith('.css'),
			)

			const cssLinks = cssFiles
				.map(file => `<link rel="stylesheet" href="./${file}">`)
				.join('\n    ')

			const htmlContent = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>App</title>
    ${cssLinks}
</head>
<body>
<div id="root"></div>
<script type="module" src="./index.js"></script>
</body>
</html>`
			this.emitFile({
				type: 'asset',
				fileName: 'index.html',
				source: htmlContent,
			})
		},
	}
}

export {
	generateIndexHtml,
}
