// e2e: un navegador real usa el front desplegado, y el front le habla a la API de verdad.
// Sin mocks: nada de page.route ni de respuestas inventadas.
import { test, expect } from '@playwright/test'

// Un nombre que ninguna otra corrida usa (QA es una sola base compartida).
function nombreUnico(prefijo) {
  const hora = new Date().toISOString().replace(/\D/g, '').slice(0, 14)
  return `${prefijo}-${hora}-${Math.random().toString(36).slice(2, 6)}`
}

// Espera a que la app termine de cargar. Con el servicio dormido (cold start) la
// primera carga puede tardar ~1 minuto: se espera, no se saltea.
async function abrirApp(page) {
  await page.goto('/')
  await expect(page.getByRole('heading', { name: /^Total:/ })).toBeVisible({ timeout: 90_000 })
}

async function cargarGasto(page, descripcion, monto) {
  await page.getByLabel('Descripción', { exact: true }).fill(descripcion)
  await page.getByLabel('Monto', { exact: true }).fill(monto)
  await page.getByRole('button', { name: 'Agregar', exact: true }).click()
}

const fila = (page, descripcion) => page.getByRole('row').filter({ hasText: descripcion })

test.describe('Gastos, de punta a punta', () => {
  // Red de seguridad: si un test falla a mitad de camino, borra lo que dejo (por la API,
  // a traves del mismo front). El flujo de cada test ya limpia y lo comprueba por su cuenta.
  const usados = []
  test.afterEach(async ({ request }) => {
    const res = await request.get('/api/gastos')
    if (!res.ok()) return
    for (const g of await res.json()) {
      if (usados.includes(g.descripcion)) await request.delete(`/api/gastos/${g.id}`)
    }
    usados.length = 0
  })

  test('crear un gasto, verlo en la lista y borrarlo', async ({ page }) => {
    const descripcion = nombreUnico('e2e-alta')
    usados.push(descripcion)
    await abrirApp(page)

    await cargarGasto(page, descripcion, '123.45')
    await expect(fila(page, descripcion)).toBeVisible()
    await expect(fila(page, descripcion)).toContainText('$123.45')

    await page.getByRole('button', { name: `Eliminar ${descripcion}` }).click()
    await expect(fila(page, descripcion)).toHaveCount(0)
  })

  test('un monto invalido muestra el error y no crea nada', async ({ page, request }) => {
    const descripcion = nombreUnico('e2e-invalido')
    usados.push(descripcion)
    await abrirApp(page)

    await page.getByLabel('Descripción', { exact: true }).fill(descripcion)
    await page.getByLabel('Monto', { exact: true }).fill('0')

    await expect(page.getByRole('alert')).toContainText('El monto tiene que ser mayor a 0')
    await expect(page.getByRole('button', { name: 'Agregar', exact: true })).toBeDisabled()

    // Nada se creo: ni en la pantalla ni en la base.
    await expect(fila(page, descripcion)).toHaveCount(0)
    const lista = await (await request.get('/api/gastos')).json()
    expect(lista.filter((g) => g.descripcion === descripcion)).toHaveLength(0)
  })

  // El flujo de todos los dias: cargo un gasto, cierro la pagina, vuelvo, y sigue ahi;
  // lo marco como pagado y eso tambien sobrevive. Toca la base de punta a punta.
  test('el gasto sigue ahi al recargar y se puede marcar como pagado', async ({ page }) => {
    const descripcion = nombreUnico('e2e-diario')
    usados.push(descripcion)
    await abrirApp(page)

    await cargarGasto(page, descripcion, '50')
    await expect(fila(page, descripcion)).toBeVisible()

    await page.reload()
    await expect(page.getByRole('heading', { name: /^Total:/ })).toBeVisible()
    await expect(fila(page, descripcion)).toContainText('pendiente')

    await page.getByRole('button', { name: `Marcar pagado ${descripcion}` }).click()
    await expect(fila(page, descripcion)).toContainText('pagado')

    await page.reload()
    await expect(page.getByRole('heading', { name: /^Total:/ })).toBeVisible()
    await expect(fila(page, descripcion)).toContainText('pagado')

    await page.getByRole('button', { name: `Eliminar ${descripcion}` }).click()
    await expect(fila(page, descripcion)).toHaveCount(0)
  })
})
