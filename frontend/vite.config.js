import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// El proxy es solo para `npm run dev` (localhost:5173). En el contenedor,
// quien traduce /api hacia el backend es nginx (ver nginx.conf).
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  test: {
    environment: 'node', // la logica a probar es pura: no hace falta DOM
    include: ['src/**/*.test.js'],
    coverage: {
      provider: 'v8',
      // Donde se guarda el reporte. El pipeline lo cambia con COVERAGE_DIR.
      reportsDirectory: process.env.COVERAGE_DIR || 'coverage',
      reporter: ['text', 'json-summary', 'html'],
      // Alcance: SOLO la logica pura en src/lib. Excluidos a proposito:
      // main.jsx (arranque de React) y App.jsx (UI/render). Son codigo de
      // presentacion: se prueba con tests de componentes/E2E, no unitarios.
      include: ['src/lib/**'],
      exclude: ['src/**/*.test.js'],
      // Umbral: si cualquiera de las dos metricas queda por debajo, vitest
      // termina con exit code 1 y el job de CI falla.
      thresholds: {
        lines: 90,
        branches: 90,
      },
    },
  },
})