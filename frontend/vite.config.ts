import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Forward API to backend during dev
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
