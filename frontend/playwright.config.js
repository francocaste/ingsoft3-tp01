import { defineConfig } from '@playwright/test'

// Dos suites, dos direcciones (se pasan por variable de entorno):
//   API_BASE_URL -> e2e/api.spec.js      integracion: le habla a la API, sin navegador
//   E2E_BASE_URL -> e2e/gastos.spec.js   e2e: un navegador usa el front
// Por defecto apuntan al compose local: API en 8080, front en 3000.
export default defineConfig({
  testDir: './e2e',

  // Cold start del free tier: despues de un rato sin trafico, Render y Neon
  // tardan ~1 minuto en despertar. Se resuelve con TIEMPO, no con sleeps ni skips.
  timeout: 90_000, // tope por test
  expect: { timeout: 15_000 }, // tope por aserciones (esperan hasta que se cumplan)

  // QA es UNA sola base compartida: los tests corren de a uno.
  fullyParallel: false,
  workers: 1,

  // Un reintento en CI absorbe un hipo de red sin esconder un bug: si falla las
  // dos veces, es rojo. La traza se graba solo en ese reintento.
  retries: process.env.CI ? 1 : 0,

  reporter: [['list'], ['html', { open: 'never' }]],

  use: {
    baseURL: process.env.E2E_BASE_URL || 'http://localhost:3000',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    actionTimeout: 15_000,
    navigationTimeout: 60_000,
  },
})
