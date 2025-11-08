import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: '../web',
    emptyOutDir: true, // Don't delete existing files in web directory
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8125',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:8125',
        ws: true,
      }
    }
  }
})

