import {resolve} from 'path'
import react from '@vitejs/plugin-react'
import {defineConfig} from 'vite'
import {generateIndexHtml} from './generateIndexHtml'

export default defineConfig(({command}) => ({
  plugins: [
    react(),
    generateIndexHtml(),
  ],
  base: './',
  root: resolve('src'),
  css: {
    modules: {
      localsConvention: 'camelCase',
    },
  },
  build: {
    outDir: resolve('dist'),
    emptyOutDir: true,
    target: 'es2022',
    rollupOptions: {
      input: {
        main: resolve('src', command === 'serve'
            ? 'index-dev.ts'
            : 'index.ts'),
      },
      output: {
        entryFileNames: 'index.js',
        inlineDynamicImports: true,
      },
    },
  },
}))
